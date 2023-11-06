package ticket

import (
	"context"
	"ticket/internal/core/domain"
	"ticket/pkg/errors"
	"time"
)

func (s Service) CreateTicket(ctx context.Context, user *domain.User, data *domain.TicketRequest) (*domain.TicketResponse, error) {
	now := time.Now().UTC()
	data.Status = domain.StatusToDo
	data.Archive = false
	data.CreateAt = now
	data.UpdateAt = now
	data.User = domain.TicketUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		ImageUrl: user.ImageUrl,
	}

	ticket, err := s.ticketRepo.Insert(ctx, data)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
func (s Service) GetTicketList(ctx context.Context, query domain.GetTicketListQuery) (*domain.TicketPaginationResponse, error) {
	query.PaginationQuery.Sort = domain.QueryUpdateAt
	query.PaginationQuery.SortDirection = domain.SortDirectionDesc
	query.Archive = false

	out, err := s.ticketRepo.FindByQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s Service) GetTicketByID(ctx context.Context, id string) (*domain.TicketResponse, error) {
	out, err := s.ticketRepo.FindOneByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (s Service) UpdateTicket(ctx context.Context, user *domain.User, data *domain.UpdateTicketRequest) (*domain.TicketResponse, error) {
	old, err := s.ticketRepo.FindOneByID(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	data.UpdateAt = time.Now().UTC()
	new, err := s.ticketRepo.UpdateOneByID(ctx, data)
	if err != nil {
		return nil, err
	}

	history := &domain.TicketHistoryRequest{
		TicketID: data.ID,
		CreateAt: time.Now().UTC(),
		User: domain.TicketUser{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			ImageUrl: user.ImageUrl,
		},
		From: domain.BaseTicket{
			Name:    old.Name,
			Deatail: old.Deatail,
			Status:  old.Status,
			Archive: old.Archive,
		},
		To: domain.BaseTicket{
			Name:    new.Name,
			Deatail: new.Deatail,
			Status:  new.Status,
			Archive: new.Archive,
		},
	}
	_, err = s.ticketHistoryRepo.Insert(ctx, history)
	if err != nil {
		return nil, err
	}
	return new, nil
}
func (s Service) DeleteTicket(ctx context.Context, id string) (*domain.TicketResponse, error) {
	out, err := s.ticketRepo.DeleteOneByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = s.ticketHistoryRepo.DeleteByTicketID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = s.ticketCommentRepo.DeleteByTicketID(ctx, id)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s Service) GetTicketHistoryList(ctx context.Context, query domain.GetTicketHistoryListQuery) (*domain.TicketHistoryPaginationResponse, error) {
	query.PaginationQuery.Sort = domain.QueryUpdateAt
	query.PaginationQuery.SortDirection = domain.SortDirectionDesc
	out, err := s.ticketHistoryRepo.FindByQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s Service) CreateTicketComment(ctx context.Context, user *domain.User, data *domain.TicketCommentRequest) (*domain.TicketCommentResponse, error) {
	now := time.Now().UTC()
	data.CreateAt = now
	data.UpdateAt = now
	data.User = domain.TicketUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		ImageUrl: user.ImageUrl,
	}
	out, err := s.ticketCommentRepo.Insert(ctx, data)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s Service) GetTicketCommentList(ctx context.Context, query domain.GetTicketCommentListQuery) (*domain.TicketCommentPaginationResponse, error) {
	query.PaginationQuery.Sort = domain.QueryUpdateAt
	query.PaginationQuery.SortDirection = domain.SortDirectionDesc
	out, err := s.ticketCommentRepo.FindByQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	return out, nil

}

func (s Service) UpdateTicketComment(ctx context.Context, user *domain.User, data *domain.UpdateTicketCommentRequest) (*domain.TicketCommentResponse, error) {
	comment, err := s.ticketCommentRepo.FindOneByID(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	if user.ID != comment.User.ID {
		return nil, errors.ErrPermissionDenied
	}

	data.UpdateAt = time.Now().UTC()
	out, err := s.ticketCommentRepo.UpdateOneByID(ctx, data)
	if err != nil {
		return nil, err
	}
	return out, nil

}
func (s Service) DeleteTicketComment(ctx context.Context, user *domain.User, id string) (*domain.TicketCommentResponse, error) {
	comment, err := s.ticketCommentRepo.FindOneByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.ID != comment.User.ID {
		return nil, errors.ErrPermissionDenied
	}

	out, err := s.ticketCommentRepo.DeleteOneByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return out, nil
}
