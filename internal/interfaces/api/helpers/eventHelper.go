package helpers

import (
	"context"
	"ginapp/internal/interfaces/api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

func AddEvent(ctx context.Context, eventCollection *mongo.Collection, event models.Event) error {
	_, err := eventCollection.InsertOne(ctx, event)
	if err != nil {
		return err
	}
	return nil
}
