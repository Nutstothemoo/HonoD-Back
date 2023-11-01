package controllers

import (
	"context"
	"net/http"
	"time"

	"ginapp/database"
	"ginapp/api/models"
	"ginapp/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var EventCollection *mongo.Collection = database.OpenCollection(database.Client, "Events")

func (app *Application) GetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		events, err := database.GetEvents(ctx, app.prodCollection)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, events)

	}
		

}

func (app *Application) GetEventByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := c.Param("id")
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		event, err := database.GetEventByID(ctx, app.prodCollection, eventID)
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
		err := database.AddEvent(ctx, event)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully added the event")
	}
}

