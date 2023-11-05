package ticket

import (
	"ticket/internal/core/port"

	"go.uber.org/zap"
)

type Config struct {
	TicketRepo        port.TicketRepository
	TicketHistoryRepo port.TicketHistoryRepository
	TicketCommentRepo port.TicketCommentRepository
}

type Service struct {
	logger            *zap.SugaredLogger
	ticketRepo        port.TicketRepository
	ticketHistoryRepo port.TicketHistoryRepository
	ticketCommentRepo port.TicketCommentRepository
}

func New(logger *zap.SugaredLogger, cfg Config) port.TicketService {
	return &Service{
		logger:            logger,
		ticketRepo:        cfg.TicketRepo,
		ticketHistoryRepo: cfg.TicketHistoryRepo,
		ticketCommentRepo: cfg.TicketCommentRepo,
	}
}
