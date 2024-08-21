package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Artist struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	Name         string             `json:"name"`
	Genre        string             `json:"genre"`
	Description  string             `json:"description"`
	Country      string             `json:"country"`
	SoundCloudID string             `json:"soundcloudID"`
	Artworks     []Artwork          `json:"artworks"`
}
