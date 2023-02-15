package budget

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Budget struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name,omitempty" validate:"required"`
	Materials []BudgetMaterial   `bson:"materials" json:"materials,omitempty" validate:"required"`
	Price     float64            `bson:"price" json:"price,omitempty" validate:"required"`
}

type BudgetMaterial struct {
	ID        primitive.ObjectID      `bson:"_id" json:"id,omitempty"`
	Name      string                  `bson:"name" json:"name,omitempty"`
	Dimension BudgetMaterialDimension `bson:"dimension" json:"dimension" validate:"required"`
	Quantity  float32                 `bson:"quantity" json:"quantity" validate:"required"`
	Price     float64                 `bson:"price" json:"price" validate:"required"`
}

type BudgetMaterialDimension struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" validate:"required"`
	Metric   string             `bson:"metric" json:"metric"`
	Quantity float64            `bson:"quantity" json:"quantity"`
	Price    float64            `bson:"price" json:"price"`
}
