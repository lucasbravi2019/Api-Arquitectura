package budget

import "go.mongodb.org/mongo-driver/bson/primitive"

type BudgetDTO struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `json:"name"`
	Materials []MaterialsDTO     `json:"materials"`
	Price     float64            `json:"price"`
}

type MaterialsDTO struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `json:"name"`
	Price     float64            `json:"price"`
	Dimension DimensionDTO       `json:"dimension"`
	Quantity  float64            `json:"quantity"`
}

type DimensionDTO struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Metric   string             `bson:"metric,omitempty" json:"metric,omitempty"`
	Quantity float64            `bson:"quantity,omitempty" json:"quantity,omitempty"`
	Price    float64            `bson:"price,omitempty" json:"price,omitempty"`
}

type BudgetNameDTO struct {
	Name string `json:"name" validate:"required"`
}
