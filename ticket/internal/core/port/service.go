package port

import (
	"context"
	"ticket/internal/core/domain"
)

type AuthService interface {
	Login(ctx context.Context, data domain.LoginRequest) (*domain.TokenResponse, error)
	Logout(ctx context.Context, data domain.LogoutRequest) error
	RefreshAccessToken(ctx context.Context, data domain.RefreshTokenRequest) (*domain.TokenResponse, error)
	Authentication(ctx context.Context, accessToken string) (*domain.AuthenticationResponse, error)
	Authorization(ctx context.Context, role, method, path string) (*domain.AuthorizationResponse, error)
}

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type TicketService interface {
	CreateTicket(ctx context.Context, data *domain.TicketRequest) (*domain.TicketResponse, error)
	GetTicketList(ctx context.Context, query domain.GetTicketListQuery) (*domain.TicketPaginationResponse, error)
	GetTicketByID(ctx context.Context, id string) (*domain.TicketResponse, error)
	UpdateTicket(ctx context.Context, data *domain.UpdateTicketRequest) (*domain.TicketResponse, error)
	DeleteTicket(ctx context.Context, id string) (*domain.TicketResponse, error)
	GetTicketHistoryList(ctx context.Context, query domain.GetTicketHistoryListQuery) (*domain.TicketHistoryPaginationResponse, error)
	CreateTicketComment(ctx context.Context, data *domain.TicketCommentRequest) (*domain.TicketCommentResponse, error)
	GetTicketCommentList(ctx context.Context, query domain.GetTicketCommentListQuery) (*domain.TicketCommentPaginationResponse, error)
	UpdateTicketComment(ctx context.Context, data *domain.UpdateTicketCommentRequest) (*domain.TicketCommentResponse, error)
	DeleteTicketComment(ctx context.Context, id string) (*domain.TicketCommentResponse, error)
}
