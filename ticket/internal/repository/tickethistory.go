package repository

import (
	"context"
	"math"
	"ticket/internal/core/domain"
	"ticket/internal/core/port"
	"ticket/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

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
		filter["ticket_id"] = toObjectID(query.TicketID)
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

func (repo *ticketHistoryRepository) DeleteByTicketID(ctx context.Context, id string) (int64, error) {
	query := bson.M{
		"ticket_id": toObjectID(id),
	}
	count, err := repo.deleteMany(ctx, query, options.Delete())
	if err != nil {
		return 0, err
	}

	return count, nil
}
