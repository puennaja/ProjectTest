package repository

import (
	"time"

	"ticket/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// repository
type userRepository struct {
	repository[user]
}
type ticketRepository struct {
	repository[ticket]
}
type ticketHistoryRepository struct {
	repository[ticketHistory]
}
type ticketCommentRepository struct {
	repository[ticketComment]
}

// struct
type user struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	ImageUrl string             `bson:"image_url"`
	Password string             `bson:"pwd"`
	Role     string             `bson:"role"`
	CreateAt time.Time          `bson:"create_at"`
	UpdateAt time.Time          `bson:"update_at"`
}

func (u *user) toDomain() *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{
		ID:       u.ID.Hex(),
		Name:     u.Name,
		Email:    u.Email,
		ImageUrl: u.ImageUrl,
		Password: u.Password,
		Role:     u.Role,
		CreateAt: u.CreateAt,
		UpdateAt: u.UpdateAt,
	}
}

type ticketUser struct {
	ID       primitive.ObjectID `bson:"id"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	ImageUrl string             `bson:"image_url"`
}

func (b ticketUser) fromDomain(d domain.TicketUser) ticketUser {
	return ticketUser{
		ID:       toObjectID(d.ID),
		Name:     d.Name,
		Email:    d.Email,
		ImageUrl: d.ImageUrl,
	}
}

func (b ticketUser) toDomain() domain.TicketUser {
	return domain.TicketUser{
		ID:       b.ID.Hex(),
		Name:     b.Name,
		Email:    b.Email,
		ImageUrl: b.ImageUrl,
	}
}

type baseTicket struct {
	Name    string `bson:"name"`
	Deatail string `bson:"deatail"`
	Status  string `bson:"status"`
	Archive bool   `bson:"archive"`
}

func (b baseTicket) fromDomain(d domain.BaseTicket) baseTicket {
	b.Name = d.Name
	b.Deatail = d.Deatail
	b.Status = d.Status
	b.Archive = d.Archive
	return b
}

func (b baseTicket) toDomain() domain.BaseTicket {
	return domain.BaseTicket{
		Name:    b.Name,
		Deatail: b.Deatail,
		Status:  b.Status,
		Archive: b.Archive,
	}
}

type ticket struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Deatail  string             `bson:"deatail"`
	Status   string             `bson:"status"`
	Archive  bool               `bson:"archive"`
	CreateAt time.Time          `bson:"create_at"`
	UpdateAt time.Time          `bson:"update_at"`
	User     ticketUser         `bson:"user"`
}

func (t *ticket) fromDomain(d *domain.TicketRequest) *ticket {
	t.Name = d.Name
	t.Deatail = d.Deatail
	t.Status = d.Status
	t.Archive = d.Archive
	t.CreateAt = d.CreateAt
	t.UpdateAt = d.UpdateAt
	t.User = ticketUser{}.fromDomain(d.User)
	return t
}

func (t *ticket) toDomain() *domain.TicketResponse {
	if t == nil {
		return nil
	}
	return &domain.TicketResponse{
		ID:       t.ID.Hex(),
		Name:     t.Name,
		Deatail:  t.Deatail,
		Status:   t.Status,
		Archive:  t.Archive,
		CreateAt: t.CreateAt,
		UpdateAt: t.UpdateAt,
		User:     t.User.toDomain(),
	}
}

type ticketList []ticket

func (tl ticketList) toDomains() domain.TicketResponseList {
	out := make(domain.TicketResponseList, 0)
	for _, t := range tl {
		out = append(out, *t.toDomain())
	}
	return out
}

type updateTicket struct {
	Name     string    `bson:"name,omitempty"`
	Deatail  string    `bson:"deatail,omitempty"`
	Status   string    `bson:"status,omitempty"`
	Archive  *bool     `bson:"archive,omitempty"`
	UpdateAt time.Time `bson:"update_at,omitempty"`
}

func (t *updateTicket) fromDomain(d *domain.UpdateTicketRequest) *updateTicket {
	t.Name = d.Name
	t.Deatail = d.Deatail
	t.Status = d.Status
	t.Archive = d.Archive
	t.UpdateAt = d.UpdateAt
	return t
}

type ticketHistory struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	TicketID primitive.ObjectID `bson:"ticket_id"`
	User     ticketUser         `bson:"user"`
	CreateAt time.Time          `bson:"create_at"`
	From     baseTicket         `bson:"from"`
	To       baseTicket         `bson:"to"`
}

func (t *ticketHistory) fromDomain(d *domain.TicketHistoryRequest) *ticketHistory {
	t.TicketID = toObjectID(d.TicketID)
	t.User = ticketUser{}.fromDomain(d.User)
	t.CreateAt = d.CreateAt
	t.From = baseTicket{}.fromDomain(d.From)
	t.To = baseTicket{}.fromDomain(d.To)
	return t
}

func (t *ticketHistory) toDomain() *domain.TicketHistoryResponse {
	if t == nil {
		return nil
	}
	return &domain.TicketHistoryResponse{
		ID:       t.ID.Hex(),
		TicketID: t.TicketID.Hex(),
		User:     t.User.toDomain(),
		CreateAt: t.CreateAt,
		From:     t.From.toDomain(),
		To:       t.To.toDomain(),
	}
}

type ticketHistoryList []ticketHistory

func (tl ticketHistoryList) toDomains() domain.TicketHistoryResponseList {
	out := make(domain.TicketHistoryResponseList, 0)
	for _, t := range tl {
		out = append(out, *t.toDomain())
	}
	return out
}

type ticketComment struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	TicketID primitive.ObjectID `bson:"ticket_id"`
	User     ticketUser         `bson:"user"`
	Comment  string             `bson:"comment"`
	CreateAt time.Time          `bson:"create_at"`
	UpdateAt time.Time          `bson:"update_at"`
}

func (t *ticketComment) fromDomain(d *domain.TicketCommentRequest) *ticketComment {
	t.TicketID = toObjectID(d.TicketID)
	t.Comment = d.Comment
	t.User = ticketUser{}.fromDomain(d.User)
	t.CreateAt = d.CreateAt
	t.UpdateAt = d.UpdateAt
	return t
}

func (t *ticketComment) toDomain() *domain.TicketCommentResponse {
	if t == nil {
		return nil
	}
	return &domain.TicketCommentResponse{
		ID:       t.ID.Hex(),
		TicketID: t.TicketID.Hex(),
		User:     t.User.toDomain(),
		Comment:  t.Comment,
		CreateAt: t.CreateAt,
		UpdateAt: t.UpdateAt,
	}
}

type ticketCommentList []ticketComment

func (tl ticketCommentList) toDomains() domain.TicketCommentResponseList {
	out := make(domain.TicketCommentResponseList, 0)
	for _, t := range tl {
		out = append(out, *t.toDomain())
	}
	return out
}

type updateTicketComment struct {
	Comment  string    `bson:"comment,omitempty"`
	UpdateAt time.Time `bson:"update_at"`
}

func (t *updateTicketComment) fromDomain(d *domain.UpdateTicketCommentRequest) *updateTicketComment {
	t.Comment = d.Comment
	t.UpdateAt = d.UpdateAt
	return t
}
