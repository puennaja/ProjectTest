package repository

import (
	"context"
	"math"
	"strings"
	"ticket/internal/core/domain"
	"ticket/internal/core/port"
	"ticket/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

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
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + strings.TrimSpace(query.Name) + ".*", Options: "i"}}
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
