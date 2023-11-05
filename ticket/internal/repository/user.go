package repository

import (
	"context"
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

type userRepository struct {
	repository[user]
}

type user struct {
	baseUser
	password string    `bson:"pwd"`
	role     string    `bson:"role"`
	createAt time.Time `bson:"create_at"`
	updateAt time.Time `bson:"update_at"`
}
type baseUser struct {
	id       primitive.ObjectID `bson:"_id"`
	name     string             `bson:"name"`
	email    string             `bson:"email"`
	imageUrl string             `bson:"image_url"`
}

func (b baseUser) fromDomain(d domain.BaseUser) baseUser {
	return baseUser{
		name:     d.Name,
		email:    d.Email,
		imageUrl: d.ImageUrl,
	}
}

func (b baseUser) toDomain() domain.BaseUser {
	return domain.BaseUser{
		ID:       b.id.Hex(),
		Name:     b.name,
		Email:    b.email,
		ImageUrl: b.imageUrl,
	}
}

func (u *user) toDomain() *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{
		BaseUser: u.baseUser.toDomain(),
		Password: u.password,
		Role:     u.role,
		CreateAt: u.createAt,
		UpdateAt: u.updateAt,
	}
}

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
