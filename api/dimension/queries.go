package dimension

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDimensionById(packageId primitive.ObjectID) bson.M {
	return bson.M{"_id": packageId}
}

func UpdateDimensionById(body Dimension) bson.M {
	return bson.M{"$set": bson.M{
		"metric":   body.Metric,
		"quantity": body.Quantity,
	}}
}
