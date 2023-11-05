package user

import (
	"context"
	"ticket/internal/core/domain"
)

func (s Service) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.FindOneByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s Service) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
