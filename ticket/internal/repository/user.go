package repository

import (
	"context"
	"ticket/internal/core/domain"
	"ticket/internal/core/port"
	"ticket/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func NewUserRepository(logger *zap.SugaredLogger, mc *mongo.Client, db string) port.UserRepository {
	var repo = &userRepository{repository[user]{
		mc:             mc,
		db:             mc.Database(db),
		col:            "user",
		errInterceptor: userErrorInterceptor,
	}}

	if err := repo.createIndex(); err != nil {
		logger.Fatal(err.Error())
	}

	return repo
}

func (repo *userRepository) createIndex() error {
	col := repo.collection()
	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"email": desc,
		},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func userErrorInterceptor(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return errors.ErrUserNotFound
	}

	return errors.ErrMongo.SetError(err)
}

func (repo *userRepository) FindOneByID(ctx context.Context, id string) (*domain.User, error) {
	query := repo.buildQueryByID(id)
	v, err := repo.findOne(ctx, query, options.FindOne())
	return v.toDomain(), err
}

func (repo *userRepository) FindOneByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := bson.M{
		"email": email,
	}
	v, err := repo.findOne(ctx, query, options.FindOne())
	return v.toDomain(), err
}
