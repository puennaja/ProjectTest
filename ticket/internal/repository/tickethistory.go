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

type ticketHistoryRepository struct {
	repository[ticketHistory]
}

type ticketHistoryList []ticketHistory
type ticketHistory struct {
	id       primitive.ObjectID `bson:"_id"`
	ticketID primitive.ObjectID `bson:"ticket_id"`
	user     baseUser           `bson:"user"`
	createAt time.Time          `bson:"create_at"`
	from     baseTicket         `bson:"from"`
	to       baseTicket         `bson:"to"`
}

func (t *ticketHistory) fromDomain(d *domain.TicketHistoryRequest) *ticketHistory {
	t.ticketID = toObjectID(d.TicketID)
	t.user = baseUser{}.fromDomain(d.User)
	t.createAt = d.CreateAt
	t.from = baseTicket{}.fromDomain(d.From)
	t.to = baseTicket{}.fromDomain(d.To)
	return t
}

func (t *ticketHistory) toDomain() *domain.TicketHistoryResponse {
	if t == nil {
		return nil
	}
	return &domain.TicketHistoryResponse{
		ID:       t.id.Hex(),
		TicketID: t.ticketID.Hex(),
		User:     t.user.toDomain(),
		CreateAt: t.createAt,
		From:     t.from.toDomain(),
		To:       t.from.toDomain(),
	}
}

func (tl ticketHistoryList) toDomains() domain.TicketHistoryResponseList {
	out := make(domain.TicketHistoryResponseList, 0)
	for _, t := range tl {
		out = append(out, *t.toDomain())
	}
	return out
}

func ticketHistoryErrorInterceptor(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.ErrTicketHistoryNotFound
	}

	return errors.ErrMongo.SetError(err)
}

func NewTicketHistoryRepository(logger *zap.SugaredLogger, mc *mongo.Client, db string) port.TicketHistoryRepository {
	var repo = &ticketHistoryRepository{repository[ticketHistory]{
		mc:             mc,
		db:             mc.Database(db),
		col:            "ticket_history",
		errInterceptor: ticketHistoryErrorInterceptor,
	}}

	if err := repo.createIndex(); err != nil {
		logger.Fatal(err.Error())
	}

	return repo
}

func (repo *ticketHistoryRepository) createIndex() error {
	col := repo.collection()
	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"ticket_id": desc,
		},
		Options: options.Index(),
	})
	return err
}

func (repo *ticketHistoryRepository) Insert(ctx context.Context, data *domain.TicketHistoryRequest) (*domain.TicketHistoryResponse, error) {
	v := new(ticketHistory).fromDomain(data)
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

func (repo *ticketHistoryRepository) FindByQuery(ctx context.Context, query domain.GetTicketHistoryListQuery) (*domain.TicketHistoryPaginationResponse, error) {
	filter := bson.M{}
	if query.TicketID != "" {
		filter["ticket_id"] = query.TicketID
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

	return &domain.TicketHistoryPaginationResponse{
		PaginationQuery: query.PaginationQuery,
		TotalRows:       total,
		TotalPages:      int(math.Ceil(float64(total) / float64(query.PaginationQuery.Limit))),
		Rows:            ticketHistoryList(res).toDomains(),
	}, err
}
