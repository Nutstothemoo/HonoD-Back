package controllers

import (
	"context"
	"ginapp/api/models"
	"ginapp/database"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := EventCollection.FindOneAndUpdate(ctx, models.Event{ID: event.ID}, event)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully updated the event")
	}
}

func DeleteEvent() gin.HandlerFunc {
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