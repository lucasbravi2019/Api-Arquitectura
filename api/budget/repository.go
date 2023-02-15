package budget

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	db *mongo.Collection
}

type BudgetRepository interface {
	FindAllBudgets() *[]BudgetDTO
	FindBudgetByOID(oid *primitive.ObjectID) *BudgetDTO
	FindBudgetsByDimensionId(oid *primitive.ObjectID) []BudgetDTO
	CreateBudget(budget *BudgetNameDTO) *primitive.ObjectID
	UpdateBudgetName(oid *primitive.ObjectID, budgetName *BudgetNameDTO) error
	AddMaterialToBudget(oid *primitive.ObjectID, recipe *BudgetMaterial) error
	RemoveMaterialFromBudget(oid *primitive.ObjectID, budget *BudgetMaterial) error
	DeleteBudget(oid *primitive.ObjectID) error
	RemoveMaterialByDimensionId(packageId *primitive.ObjectID) error
	UpdateBudgetByIdPrice(recipeId *primitive.ObjectID) error
	UpdateMaterialDimensionPrice(packageId *primitive.ObjectID, price float64) error
	UpdateMaterialsPrice(packageId *primitive.ObjectID, recipe BudgetDTO) error
	UpdateBudgetsPrice() error
}

var budgetRepositoryInstance *repository

func (r *repository) FindAllBudgets() *[]BudgetDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := r.db.Find(ctx, All())

	var recipes *[]BudgetDTO = &[]BudgetDTO{}

	if err != nil {
		log.Println(err.Error())
		return recipes
	}

	err = cursor.All(ctx, recipes)

	if err != nil {
		log.Println(err.Error())
	}

	return recipes
}

func (r *repository) FindBudgetByOID(oid *primitive.ObjectID) *BudgetDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var recipe *BudgetDTO = &BudgetDTO{}

	err := r.db.FindOne(ctx, GetRecipeById(*oid)).Decode(recipe)

	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return recipe
}

func (r *repository) FindBudgetsByDimensionId(dimensionId *primitive.ObjectID) []BudgetDTO {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := r.db.Find(ctx, GetBudgetByDimensionId(*dimensionId))

	var budgets []BudgetDTO = []BudgetDTO{}

	if err != nil {
		log.Println(err.Error())
		return budgets
	}

	err = cursor.All(ctx, &budgets)

	if err != nil {
		log.Println(err.Error())
	}

	return budgets
}

func (r *repository) CreateBudget(recipe *BudgetNameDTO) *primitive.ObjectID {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := r.db.InsertOne(ctx, recipe)

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

func (r *repository) UpdateBudgetName(oid *primitive.ObjectID, budgetName *BudgetNameDTO) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(ctx, GetRecipeById(*oid), UpdateRecipeName(*budgetName))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) AddMaterialToBudget(oid *primitive.ObjectID, budget *BudgetMaterial) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(ctx, GetRecipeById(*oid), AddIngredientToRecipe(*budget))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) RemoveMaterialFromBudget(oid *primitive.ObjectID, budget *BudgetMaterial) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(ctx, GetRecipeById(*oid), RemoveMaterialFromBudget(*budget))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) DeleteBudget(oid *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := r.db.DeleteOne(ctx, GetRecipeById(*oid))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) RemoveMaterialByDimensionId(dimensionId *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)
	defer cancel()

	_, err := r.db.UpdateMany(ctx, GetBudgetByDimensionId(*dimensionId), RemoveDimensionFromBudget(*dimensionId))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) UpdateBudgetByIdPrice(budgetId *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(ctx, GetRecipeById(*budgetId), SetBudgetPrice())

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) UpdateBudgetsPrice() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := r.db.UpdateMany(ctx, All(), SetBudgetPrice())

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) UpdateMaterialDimensionPrice(dimensionId *primitive.ObjectID, price float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := r.db.UpdateMany(ctx, GetBudgetByDimensionId(*dimensionId), SetMaterialDimensionPrice(price), GetArrayFiltersForMaterialsByDimensionId(*dimensionId))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (r *repository) UpdateMaterialsPrice(packageId *primitive.ObjectID, budget BudgetDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(ctx, GetBudgetByDimensionId(*packageId), SetMaterialPrice(budget))

	if err != nil {
		log.Println(err.Error())
	}

	return err
}
