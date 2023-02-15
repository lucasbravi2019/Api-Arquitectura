package dimension

import "go.mongodb.org/mongo-driver/bson/primitive"

type Dimension struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Metric   string             `bson:"metric" json:"metric" validate:"required"`
	Quantity float64            `bson:"quantity" json:"quantity" validate:"required"`
}
