package core

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertHexToObjectId(id string) *primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println("Object ID invalid")
		return nil
	}
	return &oid
}
