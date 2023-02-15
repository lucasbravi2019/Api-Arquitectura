package dimension

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	db *mongo.Collection
}

type DimensionRepository interface {
	GetDimensions() *[]Dimension
	GetDimensionById(oid *primitive.ObjectID) *Dimension
	CreateDimension(body *Dimension) *primitive.ObjectID
	UpdateDimension(oid *primitive.ObjectID, body *Dimension) error
	DeleteDimension(oid *primitive.ObjectID) error
}

var dimensionRepositoryInstance *repository

func (r *repository) GetDimensions() *[]Dimension {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)

	defer cancel()

	cursor, err := r.db.Find(ctx, bson.M{})

	var dimensions *[]Dimension = &[]Dimension{}

	if err != nil {
		log.Println(err.Error())
		return dimensions
	}

	err = cursor.All(ctx, dimensions)

	if err != nil {
		log.Println(err.Error())
	}

	return dimensions
}

func (r *repository) CreateDimension(body *Dimension) *primitive.ObjectID {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)

	defer cancel()

	result, err := r.db.InsertOne(ctx, body)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if result.InsertedID == nil {
		return nil
	}

	id := result.InsertedID.(primitive.ObjectID)

	return &id
}

func (r *repository) UpdateDimension(oid *primitive.ObjectID, body *Dimension) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)

	defer cancel()

	_, err := r.db.UpdateOne(ctx, GetDimensionById(*oid), UpdateDimensionById(*body))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) DeleteDimension(oid *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)

	defer cancel()

	_, err := r.db.DeleteOne(ctx, GetDimensionById(*oid))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) GetDimensionById(oid *primitive.ObjectID) *Dimension {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	var envase *Dimension = &Dimension{}

	err := r.db.FindOne(ctx, GetDimensionById(*oid)).Decode(envase)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return envase
}
