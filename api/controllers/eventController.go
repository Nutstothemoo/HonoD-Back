package controllers

import (
	"context"
	"ginapp/api/models"
	"ginapp/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var EventCollection *mongo.Collection = database.OpenCollection(database.Client, "Events")

func GetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		events, err := EventCollection.Find(ctx, models.Event{})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, events)

	}
}

func GetEventByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := c.Param("id")
		
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event, err := EventCollection.Find(ctx, bson.M{"_id": eventID})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, event)
	}
}

func AddEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.Event
		if err := c.BindJSON(&event); err != nil {
				c.IndentedJSON(http.StatusBadRequest, err)
				return
		}

		// Get dealerID from context
		dealerID, exists := c.MustGet("userId").(primitive.ObjectID)
		if !exists {
				c.IndentedJSON(http.StatusBadRequest, "User ID not found")
				return
		}

		// Set DealerID, CreatedAt, and UpdatedAt
		event.DealerID = dealerID              // Set DealerID to the dealerID from context
		event.CreatedAt = time.Now()           // Set CreatedAt to the current time
		event.UpdatedAt = time.Now()           // Set UpdatedAt to the current time

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := EventCollection.InsertOne(ctx, event)
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
		}
		c.IndentedJSON(201, "Successfully added the event")
}
}

func GetEventFromDateToDate() gin.HandlerFunc {
	return func (c *gin.Context) {
		fromDate := c.Param("fromDate")
		toDate := c.Param("toDate")
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.M{
			"startTime": bson.M{
					"$gte": fromDate,
					"$lte": toDate,
			},
	}
	events, err := EventCollection.Find(ctx, filter)		

	if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, events)
	}
}

func UpdateEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.Event
		if err := c.BindJSON(&event); err != nil {
				c.IndentedJSON(http.StatusBadRequest, err)
				return
		}

	
		dealerID, exists := c.MustGet("userId").(primitive.ObjectID)
		if !exists {
				c.IndentedJSON(http.StatusBadRequest, "User ID not found")
				return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var existingEvent models.Event
		err := EventCollection.FindOne(ctx, bson.M{"_id": event.ID}).Decode(&existingEvent)
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Error querying the database")
				return
		}

		if existingEvent.DealerID != dealerID {
				c.IndentedJSON(http.StatusUnauthorized, "You do not have permission to update this event")
				return
		}

		event.UpdatedAt = time.Now()

		_, err = EventCollection.UpdateOne(ctx, bson.M{"_id": event.ID}, bson.M{"$set": event})
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
		}

		c.IndentedJSON(201, "Successfully updated the event")
}
}

func DeleteEvent() gin.HandlerFunc {
	return func(c *gin.Context) {

		eventID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
				c.IndentedJSON(http.StatusBadRequest, "Invalid event ID")
				return
		}

		dealerID, exists := c.MustGet("userId").(primitive.ObjectID)
		if !exists {
				c.IndentedJSON(http.StatusBadRequest, "User ID not found")
				return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()


		var event models.Event
		err = EventCollection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&event)
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Error querying the database")
				return
		}

		if event.DealerID != dealerID {
				c.IndentedJSON(http.StatusUnauthorized, "You do not have permission to delete this event")
				return
		}

		_, err = EventCollection.DeleteOne(ctx, bson.M{"_id": eventID})
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
		}

		c.IndentedJSON(201, "Successfully deleted the event")
	}
}

func MainGetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		events, err := EventCollection.Find(ctx, models.Event{})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, events)

	}
}

func MainGetEventByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := c.Param("id")
		
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event, err := EventCollection.Find(ctx, bson.M{"_id": eventID})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, event)
	}
}

func AdminAddEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.Event
		if err := c.BindJSON(&event); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event.CreatedAt = time.Now()           // Set CreatedAt to the current time
		event.UpdatedAt = time.Now() 
		_, err := EventCollection.InsertOne(ctx, event)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully added the event")
	}
}

func AdminUpdateEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.Event
		if err := c.BindJSON(&event); err != nil {
				c.IndentedJSON(http.StatusBadRequest, err)
				return
		}

		eventID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
				c.IndentedJSON(http.StatusBadRequest, "Invalid event ID")
				return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event.UpdatedAt = time.Now()
		_, err = EventCollection.UpdateOne(ctx, bson.M{"_id": eventID}, bson.M{"$set": event})
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
		}

		c.IndentedJSON(201, "Successfully updated the event")
}
}

func AdminDeleteEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := c.Param("id")
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := EventCollection.DeleteOne(ctx, bson.M{"_id": eventID})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully deleted the event")
	}
}