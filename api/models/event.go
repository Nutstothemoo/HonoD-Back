package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID                primitive.ObjectID `bson:"_id" json:"_id"`
	Slug              string             `json:"slug"`
	StartTime         int64              `json:"startTime"`
	EndTime           int64              `json:"endTime"`
	Timezone          string             `json:"timezone"`
	Description       string             `json:"description"`
	AddressVisibility string             `json:"addressVisibility"`
	GeolocationID     primitive.ObjectID `bson:"geolocation_id" json:"geolocation_id"`
	IsFestival        bool               `json:"isFestival"`
	Name              string             `json:"name"`
	FeaturedText      string             `json:"featuredText"`
	Artworks          []Artwork          `json:"artworks"`
	CancelledAt       *int64             `json:"cancelledAt"`
	Currency          string             `json:"currency"`
	Tags              []Tag              `json:"tags"`
	DealerID          primitive.ObjectID `bson:"dealer_id" json:"dealer_id"`
	LaunchedAt        int64              `json:"launchedAt"`
	IsSoldOut         bool               `json:"isSoldOut"`
	MinTicketPrice    float64            `json:"minTicketPrice"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

type Artwork struct {
	ID           string `json:"id"`
	Role         string `json:"role"`
	OriginalUrl  string `json:"originalUrl"`
	BlurDataUrl  string `json:"blurDataUrl"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Typo string `json:"typo"`
}

type Dealer struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Slug                  string `json:"slug"`
	Website               string `json:"website"`
	Logo                  string `json:"logo"`
	LogoBlurDataUrl       string `json:"logoBlurDataUrl"`
	TotalFollowing        int    `json:"totalFollowing"`
	TotalNumberOfFutureEvents int `json:"totalNumberOfFutureEvents"`
	GoogleTagManagerId    string `json:"googleTagManagerId"`
	FacebookPixelId       string `json:"facebookPixelId"`
	FacebookAccessToken    string `json:"facebookAccessToken"`
}

