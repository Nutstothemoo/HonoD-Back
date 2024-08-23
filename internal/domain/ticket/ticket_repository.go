package repositories

import (
	"context"
	models "ginapp/internal/domain/ticket/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketRepository interface {
	SearchTicket(ctx context.Context) ([]models.Ticket, error)
	SearchTicketByQuery(ctx context.Context, query string) ([]models.Ticket, error)
	GetTickets(ctx context.Context, id primitive.ObjectID) (models.Ticket, error)
	AddTicket(ctx context.Context, ticket models.Ticket) error
	UpdateTicket(ctx context.Context, id primitive.ObjectID, ticket models.Ticket) error
	DeleteTicket(ctx context.Context, id primitive.ObjectID) error
}
