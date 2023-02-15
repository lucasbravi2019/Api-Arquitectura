package dimension

import "go.mongodb.org/mongo-driver/bson/primitive"

type DimensionDTO struct {
	ID primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
}
