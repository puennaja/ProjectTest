package repository

import (
	"context"
	"math"
	"ticket/internal/core/domain"
	"ticket/internal/core/port"
	"ticket/pkg/errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type ticketRepository struct {
	repository[ticket]
}

type baseTicket struct {
	name    string `bson:"name"`
	deatail string `bson:"deatail"`
	status  string `bson:"status"`
	archive bool   `bson:"archive,omitempty"`
}
type ticketList []ticket
type ticket struct {
	baseTicket
	id       primitive.ObjectID `bson:"_id"`
	createAt time.Time          `bson:"create_at"`
	updateAt time.Time          `bson:"update_at"`
	user     baseUser           `bson:"user"`
}

func (b baseTicket) fromDomain(d domain.BaseTicket) baseTicket {
	b.name = d.Name
	b.deatail = d.Deatail
	b.status = d.Status
	b.archive = d.Archive
	return b
}

func (b baseTicket) toDomain() domain.BaseTicket {
	return domain.BaseTicket{
		Name:    b.name,
		Deatail: b.deatail,
		Status:  b.status,
		Archive: b.archive,
	}
}

func (t *ticket) fromDomain(d *domain.TicketRequest) *ticket {
	t.name = d.Name
	t.deatail = d.Deatail
	t.status = d.Status
	t.archive = d.Archive
	t.createAt = d.CreateAt
	t.updateAt = d.UpdateAt
	t.user = baseUser{}.fromDomain(d.User)
	return t
}

func (t *ticket) toDomain() *domain.TicketResponse {
	if t == nil {
		return nil
	}
	return &domain.TicketResponse{
		ID:         t.id.Hex(),
		BaseTicket: t.baseTicket.toDomain(),
		CreateAt:   t.createAt,
		UpdateAt:   t.updateAt,
		User:       t.user.toDomain(),
	}
}

func (tl ticketList) toDomains() domain.TicketResponseList {
	out := make(domain.TicketResponseList, 0)
	for _, t := range tl {
		out = append(out, *t.toDomain())
	}
	return out
}

type updateTicket struct {
	name     string    `bson:"name,omitempty"`
	deatail  string    `bson:"deatail,omitempty"`
	status   string    `bson:"status,omitempty"`
	archive  bool      `bson:"archive,omitempty"`
	updateAt time.Time `bson:"update_at,omitempty"`
}

func (t *updateTicket) fromDomain(d *domain.UpdateTicketRequest) *updateTicket {
	t.name = d.Name
	t.deatail = d.Deatail
	t.status = d.Status
	t.archive = d.Archive
	t.updateAt = d.UpdateAt
	return t
}

func ticketErrorInterceptor(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.ErrTicketNotFound
	}

	return errors.ErrMongo.SetError(err)
}

func NewTicketRepository(logger *zap.SugaredLogger, mc *mongo.Client, db string) port.TicketRepository {
	var repo = &ticketRepository{repository[ticket]{
		mc:             mc,
		db:             mc.Database(db),
		col:            "ticket",
		errInterceptor: ticketErrorInterceptor,
	}}
	return repo
}

func (repo *ticketRepository) Insert(ctx context.Context, data *domain.TicketRequest) (*domain.TicketResponse, error) {
	v := new(ticket).fromDomain(data)
	oid, err := repo.insertOne(ctx, v)
	if err != nil {
		return nil, err
	}

	query := repo.buildQueryByID(oid.Hex())
	out, err := repo.findOne(ctx, query, options.FindOne())
	if err != nil {
		return nil, err
	}
	return out.toDomain(), nil
}

func (repo *ticketRepository) FindByQuery(ctx context.Context, query domain.GetTicketListQuery) (*domain.TicketPaginationResponse, error) {
	filter := bson.M{}
	filter["archive"] = query.Archive
	if query.Name != "" {
		filter["name"] = query.Name
	}

	total, err := repo.countDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := repo.paginationFindOptions(query.PaginationQuery)
	res, err := repo.find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	return &domain.TicketPaginationResponse{
		PaginationQuery: query.PaginationQuery,
		TotalRows:       total,
		TotalPages:      int(math.Ceil(float64(total) / float64(query.PaginationQuery.Limit))),
		Rows:            ticketList(res).toDomains(),
	}, err
}

func (repo *ticketRepository) FindOneByID(ctx context.Context, id string) (*domain.TicketResponse, error) {
	query := repo.buildQueryByID(id)
	v, err := repo.findOne(ctx, query, options.FindOne())
	return v.toDomain(), err
}

func (repo *ticketRepository) UpdateOneByID(ctx context.Context, data *domain.UpdateTicketRequest) (*domain.TicketResponse, error) {
	query := repo.buildQueryByID(data.ID)
	updated := new(updateTicket).fromDomain(data)
	v, err := repo.updateOne(ctx, query, updated, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if err != nil {
		return nil, err
	}
	return v.toDomain(), nil
}

func (repo *ticketRepository) DeleteOneByID(ctx context.Context, id string) (*domain.TicketResponse, error) {
	query := repo.buildQueryByID(id)
	v, err := repo.deleteOne(ctx, query, options.FindOneAndDelete())
	return v.toDomain(), err
}
