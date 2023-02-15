package material

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	materialCollection  *mongo.Collection
	dimensionCollection *mongo.Collection
	budgetCollection    *mongo.Collection
}

type MaterialRepository interface {
	GetAllMaterials() []MaterialDTO
	FindMaterialByOID(oid *primitive.ObjectID) *MaterialDTO
	FindMaterialByPackageId(packageId *primitive.ObjectID) *MaterialDTO
	ValidateExistingMaterial(MaterialName *MaterialNameDTO) error
	CreateMaterial(Material *Material) *primitive.ObjectID
	UpdateMaterial(oid *primitive.ObjectID, dto *MaterialNameDTO) error
	DeleteMaterial(oid *primitive.ObjectID) error
	AddDimensionToMaterial(MaterialOid *primitive.ObjectID, packageOid *primitive.ObjectID, envase *MaterialDimension) error
	RemoveDimensionFromMaterials(dto MaterialDimensionDTO) error
	ChangeMaterialPrice(packageOid *primitive.ObjectID, priceDTO *MaterialDimensionPriceDTO) error
}

var materialRepositoryInstance *repository

func (r *repository) GetAllMaterials() []MaterialDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := r.materialCollection.Find(ctx, All())

	if err != nil {
		log.Println(err.Error())
	}

	var materials *[]MaterialDTO = &[]MaterialDTO{}

	err = results.All(ctx, materials)

	if err != nil {
		log.Println(err.Error())
	}

	if len(*materials) < 1 {
		return []MaterialDTO{}
	}

	return *materials
}

func (r *repository) FindMaterialByOID(oid *primitive.ObjectID) *MaterialDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var material *MaterialDTO = &MaterialDTO{}

	err := r.materialCollection.FindOne(ctx, GetMaterialById(*oid)).Decode(material)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return material
}

func (r *repository) FindMaterialByPackageId(packageId *primitive.ObjectID) *MaterialDTO {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	var material *MaterialDTO = &MaterialDTO{}

	err := r.materialCollection.FindOne(ctx, GetMaterialByDimensionId(*packageId)).Decode(material)

	if err != nil {
		return nil
	}

	return material
}

func (r *repository) CreateMaterial(material *Material) *primitive.ObjectID {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	insertResult, err := r.materialCollection.InsertOne(ctx, *material)

	if err != nil {
		log.Println(err.Error())
	}

	id := insertResult.InsertedID.(primitive.ObjectID)

	return &id
}

func (r *repository) UpdateMaterial(oid *primitive.ObjectID, dto *MaterialNameDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := r.materialCollection.UpdateOne(ctx, GetMaterialById(*oid), UpdateMaterialName(*dto))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) DeleteMaterial(oid *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := r.materialCollection.DeleteOne(ctx, GetMaterialById(*oid))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) AddDimensionToMaterial(MaterialOid *primitive.ObjectID, packageOid *primitive.ObjectID, dimension *MaterialDimension) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	_, err := r.materialCollection.UpdateOne(ctx, GetMaterialWithoutExistingDimension(*MaterialOid, *packageOid), PushDimensionIntoMaterial(*dimension))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) RemoveDimensionFromMaterials(dto MaterialDimensionDTO) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	_, err := r.materialCollection.UpdateMany(ctx, GetMaterialByDimensionId(dto.DimensionOid), PullDimensionFromMaterials(dto))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) ChangeMaterialPrice(dimensionId *primitive.ObjectID, priceDTO *MaterialDimensionPriceDTO) error {

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	_, err := r.materialCollection.UpdateOne(ctx, GetMaterialByDimensionId(*dimensionId), SetMaterialPrice(priceDTO.Price), GetArrayFilterForPackageId(*dimensionId))

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (r *repository) ValidateExistingMaterial(MaterialName *MaterialNameDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := r.materialCollection.Aggregate(ctx, GetAggregateCreateMaterials(MaterialName))

	if err != nil {
		log.Println(err.Error())
		return err
	}

	var MaterialsDuplicated *[]MaterialDTO = &[]MaterialDTO{}

	err = cursor.All(ctx, MaterialsDuplicated)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	if len(*MaterialsDuplicated) > 0 {
		return err
	}

	return nil
}
