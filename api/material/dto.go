package material

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MaterialNameDTO struct {
	Name string `json:"name" validate:"required"`
}

type MaterialDTO struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name,omitempty" json:"name,omitempty" validate:"required"`
	Dimensions []DimensionDTO     `bson:"dimensions,omitempty" json:"dimensions,omitempty"`
}

type DimensionDTO struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Metric   string             `bson:"metric,omitempty" json:"metric,omitempty"`
	Quantity float64            `bson:"quantity,omitempty" json:"quantity,omitempty"`
	Price    float64            `bson:"price,omitempty" json:"price,omitempty"`
}

type MaterialDimensionDTO struct {
	MaterialOid  primitive.ObjectID
	DimensionOid primitive.ObjectID
	Price        float64
}

type MaterialDimensionPriceDTO struct {
	Price float64 `json:"price" validate:"required"`
}

type BudgetMaterialDTO struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `json:"name"`
	Price     float64            `json:"price"`
	Dimension DimensionDTO       `json:"dimension"`
	Quantity  float64            `json:"quantity"`
}

type MaterialDetailsDTO struct {
	Metric   string  `json:"metric,omitempty"`
	Quantity float32 `json:"quantity,omitempty"`
}
