package httphdl

import (
	"ticket/internal/core/port"
	"ticket/pkg/validator"
)

type Config struct {
	Validator     validator.Validator
	AuthService   port.AuthService
	UserService   port.UserService
	TicketService port.TicketService
}

type HTTPHandler struct {
	validator     validator.Validator
	authService   port.AuthService
	userService   port.UserService
	ticketService port.TicketService
}

func NewHTTP(cfg Config) *HTTPHandler {
	return &HTTPHandler{
		validator:     cfg.Validator,
		authService:   cfg.AuthService,
		userService:   cfg.UserService,
		ticketService: cfg.TicketService,
	}
}
