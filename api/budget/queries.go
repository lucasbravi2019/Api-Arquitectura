package budget

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func All() bson.M {
	return bson.M{}
}

func GetRecipeById(oid primitive.ObjectID) bson.M {
	return bson.M{"_id": oid}
}

func UpdateRecipeName(dto BudgetNameDTO) bson.M {
	return bson.M{"$set": bson.M{"name": dto.Name}}
}

func AddIngredientToRecipe(budget BudgetMaterial) bson.M {
	return bson.M{"$addToSet": bson.M{"ingredients": budget}}
}

func RemoveMaterialFromBudget(budget BudgetMaterial) bson.M {
	return bson.M{"$pull": bson.M{"materials._id": budget.ID}}
}

func SetBudgetPrice() bson.A {
	return bson.A{bson.D{{"$set", bson.D{{"price", bson.D{{"$multiply", bson.A{bson.D{{"$sum", "$materials.price"}}, 3}}}}}}}}
}

func SetMaterialDimensionPrice(price float64) bson.D {
	return bson.D{{"$set", bson.D{{"materials.$[material].dimension.price", price}}}}
}

func SetMaterialPrice(budget BudgetDTO) bson.D {
	return bson.D{{"$set", budget}}
}

func GetArrayFiltersForMaterialsByDimensionId(dimensionId primitive.ObjectID) *options.UpdateOptions {
	return options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"materials.dimension._id": dimensionId},
		},
	})
}

func GetBudgetByDimensionId(DimensionId primitive.ObjectID) bson.M {
	return bson.M{"materials.dimension._id": DimensionId}
}

func GetBudgetByMaterialId(materialId primitive.ObjectID) bson.M {
	return bson.M{"materials._id": materialId}
}

func RemoveDimensionFromBudget(dimensionId primitive.ObjectID) bson.M {
	return bson.M{"$pull": bson.M{"ingredients": bson.M{"package._id": dimensionId}}}
}
