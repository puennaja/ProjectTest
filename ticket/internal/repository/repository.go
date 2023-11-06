package repository

import (
	"context"

	"ticket/internal/core/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	errInterceptor       func(error) error
	entity               any
	repository[T entity] struct {
		mc             *mongo.Client
		db             *mongo.Database
		col            string
		errInterceptor errInterceptor
	}
)

var (
	desc          = -1
	asc           = 1
	SortDirection = map[string]int{
		domain.SortDirectionDesc: desc,
		domain.SortDirectionAsc:  asc,
	}
)

func (repo *repository[T]) collection(opts ...*options.CollectionOptions) *mongo.Collection {
	return repo.db.Collection(repo.col, opts...)
}

func (repo *repository[T]) errorInterceptor(err error) error {
	if repo.errInterceptor != nil {
		return repo.errInterceptor(err)
	}
	return err
}

func (repo *repository[T]) countDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (out int64, err error) {
	col := repo.collection()
	out, err = col.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return
	}
	return
}

func (repo *repository[T]) find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (out []T, err error) {
	col := repo.collection()
	cursor, err := col.Find(ctx, filter, opts...)
	if err != nil {
		return
	}

	err = repo.errorInterceptor(cursor.All(ctx, &out))
	return
}

func (repo *repository[T]) findOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (out *T, err error) {
	col := repo.collection()
	err = repo.errorInterceptor(col.FindOne(ctx, filter, opts...).Decode(&out))
	return
}

func (repo *repository[T]) insertOne(ctx context.Context, in *T, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
	col := repo.collection()

	out, err := col.InsertOne(ctx, &in, opts...)
	if err != nil {
		return primitive.NilObjectID, repo.errorInterceptor(err)
	}

	return asObjectID(out.InsertedID), nil
}

func (repo *repository[T]) updateOne(ctx context.Context, query, set any, opts ...*options.FindOneAndUpdateOptions) (out *T, err error) {
	col := repo.collection()

	update := bson.M{
		"$set": set,
	}
	err = repo.errorInterceptor(col.FindOneAndUpdate(ctx, query, &update, opts...).Decode(&out))
	return
}

func (repo *repository[T]) buildQueryByID(id string) bson.M {
	oid := toObjectID(id)
	if oid == primitive.NilObjectID {
		return bson.M{
			"_id": "",
		}
	}
	return bson.M{
		"_id": oid,
	}
}

func (repo *repository[T]) paginationFindOptions(query domain.PaginationQuery) *options.FindOptions {
	skip := (query.Page - 1) * query.Limit
	opts := options.Find()
	opts.SetLimit(query.Limit)
	opts.SetSkip(skip)
	opts.SetSort(bson.D{{Key: query.Sort, Value: SortDirection[query.SortDirection]}})
	return opts
}

func (repo *repository[T]) deleteOne(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (out *T, err error) {
	col := repo.collection()
	err = repo.errorInterceptor(col.FindOneAndDelete(ctx, filter, opts...).Decode(&out))
	return
}

func (repo *repository[T]) deleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (out int64, err error) {
	col := repo.collection()
	result, err := col.DeleteMany(ctx, filter, opts...)
	err = repo.errorInterceptor(err)
	out = result.DeletedCount
	return
}
