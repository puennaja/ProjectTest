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

type ticketCommentRepository struct {
	repository[ticketComment]
}

type ticketCommentList []ticketComment
type ticketComment struct {
	id       primitive.ObjectID `bson:"_id"`
	ticketID primitive.ObjectID `bson:"ticket_id"`
	user     baseUser           `bson:"user"`
	comment  string             `bson:"comment"`
	createAt time.Time          `bson:"create_at"`
	updateAt time.Time          `bson:"update_at"`
}

func (t *ticketComment) fromDomain(d *domain.TicketCommentRequest) *ticketComment {
	t.ticketID = toObjectID(d.TicketID)
	t.comment = d.Comment
	t.user = baseUser{}.fromDomain(d.User)
	t.createAt = d.CreateAt
	t.updateAt = d.UpdateAt
	return t
}

func (t *ticketComment) toDomain() *domain.TicketCommentResponse {
	if t == nil {
		return nil
	}
	return &domain.TicketCommentResponse{
		ID:       t.id.Hex(),
		TicketID: t.ticketID.Hex(),
		User:     t.user.toDomain(),
		Comment:  t.comment,
		CreateAt: t.createAt,
		UpdateAt: t.updateAt,
	}
}

func (tl ticketCommentList) toDomains() domain.TicketCommentResponseList {
	out := make(domain.TicketCommentResponseList, 0)
	for _, t := range tl {
		out = append(out, *t.toDomain())
	}
	return out
}

type updateTicketComment struct {
	comment  string    `bson:"comment,omitempty"`
	updateAt time.Time `bson:"update_at"`
}

func (t *updateTicketComment) fromDomain(d *domain.UpdateTicketCommentRequest) *updateTicketComment {
	t.comment = d.Comment
	t.updateAt = d.UpdateAt
	return t
}

func ticketCommentErrorInterceptor(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.ErrTicketCommentNotFound
	}

	return errors.ErrMongo.SetError(err)
}

func NewTicketCommentRepository(logger *zap.SugaredLogger, mc *mongo.Client, db string) port.TicketCommentRepository {
	var repo = &ticketCommentRepository{repository[ticketComment]{
		mc:             mc,
		db:             mc.Database(db),
		col:            "ticket_comment",
		errInterceptor: ticketCommentErrorInterceptor,
	}}

	if err := repo.createIndex(); err != nil {
		logger.Fatal(err.Error())
	}
	return repo
}

func (repo *ticketCommentRepository) createIndex() error {
	col := repo.collection()
	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"ticket_id": desc,
		},
		Options: options.Index(),
	})
	return err
}

func (repo *ticketCommentRepository) Insert(ctx context.Context, data *domain.TicketCommentRequest) (*domain.TicketCommentResponse, error) {
	v := new(ticketComment).fromDomain(data)
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

func (repo *ticketCommentRepository) FindByQuery(ctx context.Context, query domain.GetTicketCommentListQuery) (*domain.TicketCommentPaginationResponse, error) {
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

	return &domain.TicketCommentPaginationResponse{
		PaginationQuery: query.PaginationQuery,
		TotalRows:       total,
		TotalPages:      int(math.Ceil(float64(total) / float64(query.PaginationQuery.Limit))),
		Rows:            ticketCommentList(res).toDomains(),
	}, err
}

func (repo *ticketCommentRepository) FindOneByID(ctx context.Context, id string) (*domain.TicketCommentResponse, error) {
	query := repo.buildQueryByID(id)
	v, err := repo.findOne(ctx, query, options.FindOne())
	return v.toDomain(), err
}

func (repo *ticketCommentRepository) UpdateOneByID(ctx context.Context, data *domain.UpdateTicketCommentRequest) (*domain.TicketCommentResponse, error) {
	query := repo.buildQueryByID(data.CommentID)
	updated := new(updateTicketComment).fromDomain(data)
	v, err := repo.updateOne(ctx, query, updated, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if err != nil {
		return nil, err
	}
	return v.toDomain(), nil
}
func (repo *ticketCommentRepository) DeleteOneByID(ctx context.Context, id string) (*domain.TicketCommentResponse, error) {
	query := repo.buildQueryByID(id)
	v, err := repo.deleteOne(ctx, query, options.FindOneAndDelete())
	return v.toDomain(), err
}
