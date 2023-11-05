package user

import (
	"ticket/internal/core/port"

	"go.uber.org/zap"
)

type Config struct {
	UserRepo port.UserRepository
}

type Service struct {
	logger   *zap.SugaredLogger
	userRepo port.UserRepository
}

func New(logger *zap.SugaredLogger, cfg Config) port.UserService {
	return &Service{
		logger:   logger,
		userRepo: cfg.UserRepo,
	}
}
