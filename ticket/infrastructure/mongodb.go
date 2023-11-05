package infrastructure

import (
	"context"
	"log"
	"ticket/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetConfig().MongoDB.Timeout)
	defer cancel()

	opts := options.Client().ApplyURI(config.GetConfig().MongoDB.URI)
	c, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c

}
