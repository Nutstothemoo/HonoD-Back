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

var ArtistCollection *mongo.Collection = database.OpenCollection(database.Client, "Artists")

func GetArtists() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		artists, err := ArtistCollection.Find(ctx, models.Artist{})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, artists)
	}
}

func GetArtistByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		artistID := c.Param("id")
		
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		artist, err := ArtistCollection.Find(ctx, bson.M{"_id": artistID})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(http.StatusOK, artist)
	}
}

func AddArtist() gin.HandlerFunc {
	return func(c *gin.Context) {
		var artist models.Artist
		if err := c.BindJSON(&artist); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := database.AddArtist(ctx, artist)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully added the artist")
	}
}

func UpdateArtist() gin.HandlerFunc {
	return func(c *gin.Context) {
		var artist models.Artist
		if err := c.BindJSON(&artist); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := ArtistCollection.FindOneAndUpdate(ctx, models.Artist{ID: artist.ID}, artist)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully updated the artist")
	}
}

func DeleteArtist() gin.HandlerFunc {
	return func(c *gin.Context) {
		artistID := c.Param("id")
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := ArtistCollection.DeleteOne(ctx, bson.M{"_id": artistID})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(201, "Successfully deleted the artist")
	}
}