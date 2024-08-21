package controllers

import (
	"context"
	"fmt"
	"ginapp/internal/infrastructure/database"
	"ginapp/internal/interfaces/api/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection *mongo.Collection = database.OpenCollection(database.Client, "Users")
var TicketCollection *mongo.Collection = database.OpenCollection(database.Client, "Tickets")
var Validate = validator.New()

// USER

func SearchTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Ticket
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := TicketCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Someting Went Wrong when finding the cursor")
			return
		}
		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productlist)
	}
}

func SearchTicketByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchproducts []models.Ticket
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := TicketCollection.Find(ctx, bson.M{"ticket_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		err = searchquerydb.All(ctx, &searchproducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchproducts)
	}
}

func GetTickets() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Invalid Event ID")
			return
		}
		var tickets []models.Ticket
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cursor, err := TicketCollection.Find(ctx, bson.M{"event_id": eventID})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get tickets")
			return
		}
		fmt.Println("cursor", cursor)
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var ticket models.Ticket
			cursor.Decode(&ticket)
			tickets = append(tickets, ticket)
		}
		if err := cursor.Err(); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Failed to get tickets")
			return
		}
		fmt.Println("tickets", tickets)
		c.IndentedJSON(200, tickets)
	}
}

// DEALER

func AddTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var ticket models.Ticket
		if err := c.BindJSON(&ticket); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}

		userIDStr, exists := c.Get("userId")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, "You are not authorized to add a ticket")
			return
		}
		userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Failed to convert userId to ObjectID")
			return
		}
		ticket.DealerID = userID
		var event models.Event
		eventIDStr := c.Param("eventId")
		eventID, err := primitive.ObjectIDFromHex(eventIDStr)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Invalid event ID")
		}
		fmt.Println(eventID)
		err = EventCollection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&event)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, "no event found with the given ID")
			return
		}
		if event.DealerID != userID {
			c.IndentedJSON(http.StatusUnauthorized, "The event does not belong to you")
			return
		}
		ticket.EventID = eventID
		ticket.CreatedAt = time.Now() // Set CreatedAt to the current time
		ticket.UpdatedAt = time.Now()
		_, err = TicketCollection.InsertOne(ctx, ticket)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully added the ticket")
	}
}

func UpdateTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticket models.Ticket
		if err := c.BindJSON(&ticket); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}

		ticketID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Invalid ticket ID")
			return
		}

		userIDStr, exists := c.Get("userId")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, "You are not authorized to add a ticket")
			return
		}
		userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Failed to convert userId to ObjectID")
			return
		}
		if ticket.DealerID != userID {
			c.IndentedJSON(http.StatusUnauthorized, "You are not authorized to update this ticket")
			return
		}

		var event models.Event
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = EventCollection.FindOne(ctx, bson.M{"_id": ticket.EventID}).Decode(&event)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		if event.DealerID != userID {
			c.IndentedJSON(http.StatusUnauthorized, "The event does not belong to you")
			return
		}

		ticket.UpdatedAt = time.Now()
		_, err = TicketCollection.UpdateOne(ctx, bson.M{"_id": ticketID}, bson.M{"$set": ticket})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(201, "Successfully updated the ticket")
	}
}

func DeleteTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		ticketID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Invalid ticket ID")
			return
		}

		var ticket models.Ticket
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = TicketCollection.FindOne(ctx, bson.M{"_id": ticketID}).Decode(&ticket)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, err)
			return
		}

		userIDStr, exists := c.Get("userId")
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, "You are not authorized to add a ticket")
			return
		}
		userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Failed to convert userId to ObjectID")
			return
		}
		if ticket.DealerID != userID {
			c.IndentedJSON(http.StatusUnauthorized, "You are not authorized to delete this ticket")
			return
		}

		_, err = TicketCollection.DeleteOne(ctx, bson.M{"_id": ticketID})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(201, "Successfully deleted the ticket")
	}
}

// ADMIN

func AdminAddTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticket models.Ticket
		if err := c.BindJSON(&ticket); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ticket.CreatedAt = time.Now() // Set CreatedAt to the current time
		ticket.UpdatedAt = time.Now()
		_, err := TicketCollection.InsertOne(ctx, ticket)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully added the ticket")
	}
}

func AdminUpdateTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ticket models.Ticket
		if err := c.BindJSON(&ticket); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}

		ticketID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Invalid ticket ID")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ticket.UpdatedAt = time.Now()
		_, err = TicketCollection.UpdateOne(ctx, bson.M{"_id": ticketID}, bson.M{"$set": ticket})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(201, "Successfully updated the ticket")
	}
}

func AdminDeleteTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		ticketID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Invalid ticket ID")
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = TicketCollection.DeleteOne(ctx, bson.M{"_id": ticketID})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully deleted the ticket")
	}
}
