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
	Geolocation       Geolocation        `json:"geolocation"`
	IsFestival        bool               `json:"isFestival"`
	Name              string             `json:"name"`
	FeaturedText      string             `json:"featuredText"`
	Artworks          []Artwork          `json:"artworks"`
	CancelledAt       *int64             `json:"cancelledAt"`
	Currency          string             `json:"currency"`
	Tags              []Tag              `json:"tags"`
	DealerID          primitive.ObjectID `bson:"dealer_id" json:"dealer_id"` // Changed from Dealer to DealerID
	LaunchedAt        int64              `json:"launchedAt"`
	IsSoldOut         bool               `json:"isSoldOut"`
	MinTicketPrice    float64            `json:"minTicketPrice"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"` // Added
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"` // Added
}

type Geolocation struct {
	Street string `json:"street"`
	Venue  string `json:"venue"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	City   City   `json:"city"`
}

type City struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location GPSPoint `json:"location"`
	Area     Area   `json:"area"`
	Country  Country `json:"country"`
	ZipCode  string `json:"zipCode"`
}

type GPSPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Area struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Country struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsoCode string `json:"isoCode"`
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

