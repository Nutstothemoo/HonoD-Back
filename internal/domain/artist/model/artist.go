package models

import (
	models "ginapp/internal/domain/event/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Artist struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	Name         string             `json:"name"`
	Genre        string             `json:"genre"`
	Description  string             `json:"description"`
	Country      string             `json:"country"`
	SoundCloudID string             `json:"soundcloudID"`
	Artworks     []models.Artwork   `json:"artworks"`
}
