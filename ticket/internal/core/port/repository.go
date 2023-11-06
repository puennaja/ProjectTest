package port

import (
	"context"

	"ticket/internal/core/domain"
)

type UserRepository interface {
	FindOneByID(ctx context.Context, id string) (*domain.User, error)
	FindOneByEmail(ctx context.Context, email string) (*domain.User, error)
}

type TicketRepository interface {
	Insert(ctx context.Context, data *domain.TicketRequest) (*domain.TicketResponse, error)
	FindByQuery(ctx context.Context, query domain.GetTicketListQuery) (*domain.TicketPaginationResponse, error)
	FindOneByID(ctx context.Context, id string) (*domain.TicketResponse, error)
	UpdateOneByID(ctx context.Context, data *domain.UpdateTicketRequest) (*domain.TicketResponse, error)
	DeleteOneByID(ctx context.Context, id string) (*domain.TicketResponse, error)
}

type TicketHistoryRepository interface {
	FindByQuery(ctx context.Context, query domain.GetTicketHistoryListQuery) (*domain.TicketHistoryPaginationResponse, error)
	Insert(ctx context.Context, data *domain.TicketHistoryRequest) (*domain.TicketHistoryResponse, error)
	DeleteByTicketID(ctx context.Context, id string) (int64, error)
}

type TicketCommentRepository interface {
	Insert(ctx context.Context, data *domain.TicketCommentRequest) (*domain.TicketCommentResponse, error)
	FindByQuery(ctx context.Context, query domain.GetTicketCommentListQuery) (*domain.TicketCommentPaginationResponse, error)
	FindOneByID(ctx context.Context, id string) (*domain.TicketCommentResponse, error)
	UpdateOneByID(ctx context.Context, data *domain.UpdateTicketCommentRequest) (*domain.TicketCommentResponse, error)
	DeleteOneByID(ctx context.Context, id string) (*domain.TicketCommentResponse, error)
	DeleteByTicketID(ctx context.Context, id string) (int64, error)
}
