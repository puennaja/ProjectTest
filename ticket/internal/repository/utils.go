package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func toObjectID(hex string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return primitive.NilObjectID
	}

	return oid
}

func asObjectID(v interface{}) primitive.ObjectID {
	oid, _ := v.(primitive.ObjectID)
	return oid
}
