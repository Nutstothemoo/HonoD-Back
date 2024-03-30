package controllers

import (
	"context"
	"fmt"
	"ginapp/api/models"
	"ginapp/database"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var EventCollection *mongo.Collection = database.OpenCollection(database.Client, "Events")

func GetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		opts := options.Find().SetSort(bson.D{{"startTime", 1}})
    events, err := EventCollection.Find(ctx, bson.D{}, opts)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		var eventsArray []bson.M
		if err = events.All(ctx, &eventsArray); err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
		}
		c.IndentedJSON(http.StatusOK, eventsArray)
	}
}

func GetEventByID() gin.HandlerFunc {
	return func(c *gin.Context) {
			var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			id, _ := primitive.ObjectIDFromHex(c.Param("id"))

			pipeline := mongo.Pipeline{
					{{"$match", bson.D{{"_id", id}}}},
					{{"$lookup", bson.D{
							{"from", "Users"},
							{"localField", "dealer_id"},
							{"foreignField", "_id"},
							{"as", "dealer"},
					}}},
					{{"$unwind", "$dealer"}},
					{{"$addFields", bson.D{
						{"geolocation", bson.D{{"$cond", bson.A{
								bson.D{{"$eq", bson.A{"$addressVisibility", "public"}}},
								"$geolocation",
								"$$REMOVE",
						}}}},
				}}},
					{{"$project", bson.D{
						{"_id", 1},
						{"addressVisibility", 1},
						{"artworks", 1},
						{"currency", 1},
						{"description", 1},
						{"endTime", 1},
						{"featuredText", 1},
						{"isFestival", 1},
						{"isSoldOut", 1},
						{"launchedAt", 1},
						{"minTicketPrice", 1},
						{"name", 1},
						{"startTime", 1},
						{"tags", 1},
						{"timezone", 1},
						{"updated_at", 1},
						{"dealer.dealerName", 1},
						{"dealer.avatar", 1},
						{"geolocation", 1},
				}}},
			}

			cursor, err := EventCollection.Aggregate(ctx, pipeline)
			if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Error finding the event"})
					return
			}

			var eventWithDealer []bson.M
			if err = cursor.All(ctx, &eventWithDealer); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding the event"})
					return
			}

			if len(eventWithDealer) == 0 {
					c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
					return
			}

			c.JSON(http.StatusOK, eventWithDealer[0])
	}
}

func AddEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.Event
		if err := c.BindJSON(&event); err != nil {
				c.IndentedJSON(http.StatusBadRequest, err)
				return
		}

		dealerIDStr, exists := c.Get("userId")
		if !exists {
				c.IndentedJSON(http.StatusUnauthorized, "You are not authorized to add a ticket")
				return
		}
		dealerID, err := primitive.ObjectIDFromHex(dealerIDStr.(string))
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Failed to convert userId to ObjectID")
				return
		}
			// If color field is empty, generate a random color
			if event.Color == "" {
				r := rand.Intn(256)
				g := rand.Intn(256)
				b := rand.Intn(256)
				event.Color = fmt.Sprintf("#%02X%02X%02X", r, g, b)
		}
		event.DealerID = dealerID              // Set DealerID to the dealerID from context
		event.CreatedAt = time.Now()           // Set CreatedAt to the current time
		event.UpdatedAt = time.Now()           // Set UpdatedAt to the current time
			// If color field is empty, generate a random color
			if event.Color == "" {
				r := rand.Intn(256)
				g := rand.Intn(256)
				b := rand.Intn(256)
				event.Color = fmt.Sprintf("#%02X%02X%02X", r, g, b)
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = EventCollection.InsertOne(ctx, event)
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
		}
		c.IndentedJSON(201, gin.H{"message": "Successfully added the event", "event": event})
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

	
		dealerIDStr, exists := c.Get("userId")
		if !exists {
				c.IndentedJSON(http.StatusUnauthorized, "You are not authorized to update a ticket")
				return
		}
		dealerID, err := primitive.ObjectIDFromHex(dealerIDStr.(string))
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Failed to convert userId to ObjectID")
				return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var existingEvent models.Event
		eventID, err := primitive.ObjectIDFromHex(c.Param("eventId"))
		if err != nil {
				c.IndentedJSON(http.StatusBadRequest, "Invalid event ID")
		}
		err = EventCollection.FindOne(ctx, bson.M{"_id": eventID}).Decode(&existingEvent)
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Error querying the database")
				return
		}

		if existingEvent.DealerID != dealerID {
				c.IndentedJSON(http.StatusUnauthorized, "You do not have permission to update this event")
				return
		}

		event.UpdatedAt = time.Now()
		_, err = EventCollection.UpdateOne(ctx, bson.M{"_id": eventID}, bson.M{"$set": event})
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

		dealerIDStr, exists := c.Get("userId")
		if !exists {
				c.IndentedJSON(http.StatusUnauthorized, "You are not authorized to add a ticket")
				return
		}
		dealerID, err := primitive.ObjectIDFromHex(dealerIDStr.(string))
		if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Failed to convert userId to ObjectID")
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