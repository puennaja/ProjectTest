package domain

import "time"

const (
	StatusToDo       string = "to_do"
	StatusInProgress string = "in_progress"
	StatusDone       string = "done"
)

type GetTicketListQuery struct {
	PaginationQuery
	Name    string `json:"name"`
	Archive bool   `json:"-"`
}
type TicketRequest struct {
	Name     string    `json:"name,omitempty" validate:"omitempty"`
	Deatail  string    `json:"deatail,omitempty" validate:"omitempty"`
	Status   string    `json:"-"`
	Archive  bool      `json:"-"`
	CreateAt time.Time `json:"-"`
	UpdateAt time.Time `json:"-"`
	User     BaseUser  `json:"-"`
}
type BaseTicket struct {
	Name    string `json:"name"`
	Deatail string `json:"deatail"`
	Status  string `json:"status"`
	Archive bool   `json:"archive,omitempty"`
}
type TicketPaginationResponse PaginationResponse[TicketResponse]
type TicketResponseList []TicketResponse
type TicketResponse struct {
	BaseTicket
	ID       string    `json:"id"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
	User     BaseUser  `json:"user"`
}
type GetTicketQuery struct {
	ID string `url:"id" validate:"required"`
}
type UpdateTicketRequest struct {
	ID       string    `url:"id" validate:"required"`
	Name     string    `json:"name,omitempty" validate:"omitempty"`
	Deatail  string    `json:"deatail,omitempty" validate:"omitempty"`
	Status   string    `json:"status,omitempty" validate:"omitempty"`
	Archive  bool      `json:"archive,omitempty" validate:"omitempty"`
	UpdateAt time.Time `json:"-"`
}
type DeleteTicketQuery struct {
	ID string `url:"id" validate:"required"`
}

type GetTicketHistoryListQuery struct {
	PaginationQuery
	TicketID string `url:"id" validate:"required"`
}
type TicketHistoryRequest struct {
	TicketID string
	User     BaseUser
	CreateAt time.Time
	From     BaseTicket
	To       BaseTicket
}
type TicketHistoryPaginationResponse PaginationResponse[TicketHistoryResponse]
type TicketHistoryResponseList []TicketHistoryResponse
type TicketHistoryResponse struct {
	ID       string     `json:"id"`
	TicketID string     `json:"ticket_id"`
	User     BaseUser   `json:"user"`
	CreateAt time.Time  `json:"create_at"`
	From     BaseTicket `json:"from"`
	To       BaseTicket `json:"to"`
}

type GetTicketCommentListQuery struct {
	PaginationQuery
	TicketID  string `url:"id" validate:"required"`
	CommentID string `url:"comment_id" validate:"required"`
}
type TicketCommentRequest struct {
	TicketID string    `url:"id" validate:"required"`
	Comment  string    `json:"comment" validate:"required"`
	User     BaseUser  `json:"-"`
	CreateAt time.Time `json:"-"`
	UpdateAt time.Time `json:"-"`
}
type UpdateTicketCommentRequest struct {
	TicketID  string    `url:"id" validate:"required"`
	CommentID string    `url:"comment_id" validate:"required"`
	Comment   string    `json:"comment" validate:"required"`
	UpdateAt  time.Time `json:"-"`
}
type DeleteTicketCommentQuery struct {
	TicketID  string `url:"id" validate:"required"`
	CommentID string `url:"comment_id" validate:"required"`
}
type TicketCommentPaginationResponse PaginationResponse[TicketCommentResponse]
type TicketCommentResponseList []TicketCommentResponse
type TicketCommentResponse struct {
	ID       string    `json:"id"`
	TicketID string    `json:"ticket_id"`
	Comment  string    `json:"comment"`
	User     BaseUser  `json:"user"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}
