package material

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lucasbravi2019/arquitectura/api/budget"
	"github.com/lucasbravi2019/arquitectura/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	materialRepository MaterialRepository
	budgetRepository   budget.BudgetRepository
}

type MaterialService interface {
	GetAllMaterials() (int, []MaterialDTO)
	CreateMaterial(r *http.Request) (int, *MaterialDTO)
	UpdateMaterial(r *http.Request) (int, *MaterialDTO)
	DeleteMaterial(r *http.Request) (int, *primitive.ObjectID)
	AddMaterialToBudget(r *http.Request) int
	ChangeMaterialPrice(r *http.Request) (int, *MaterialDTO)
}

var materialServiceInstance *service

func (s *service) GetAllMaterials() (int, []MaterialDTO) {
	Materials := s.materialRepository.GetAllMaterials()

	return http.StatusOK, Materials
}

func (s *service) CreateMaterial(r *http.Request) (int, *MaterialDTO) {
	var MaterialDto *MaterialNameDTO = &MaterialNameDTO{}

	invalidBody := core.DecodeBody(r, MaterialDto)

	if invalidBody {
		return http.StatusBadRequest, nil
	}

	var MaterialEntity *Material = &Material{
		Name:       MaterialDto.Name,
		Dimensions: []MaterialDimension{},
	}

	MaterialCreatedId := s.materialRepository.CreateMaterial(MaterialEntity)

	if MaterialCreatedId == nil {
		return http.StatusInternalServerError, nil
	}

	MaterialCreated := s.materialRepository.FindMaterialByOID(MaterialCreatedId)

	if MaterialCreated == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusCreated, MaterialCreated
}

func (s *service) UpdateMaterial(r *http.Request) (int, *MaterialDTO) {
	oid := core.ConvertHexToObjectId(mux.Vars(r)["id"])

	if oid == nil {
		return http.StatusBadRequest, nil
	}

	var Material *MaterialNameDTO = &MaterialNameDTO{}

	invalidBody := core.DecodeBody(r, Material)

	if invalidBody {
		return http.StatusBadRequest, nil
	}

	err := s.materialRepository.UpdateMaterial(oid, Material)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	MaterialUpdated := s.materialRepository.FindMaterialByOID(oid)

	if MaterialUpdated == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, MaterialUpdated
}

func (s *service) DeleteMaterial(r *http.Request) (int, *primitive.ObjectID) {
	oid := core.ConvertHexToObjectId(mux.Vars(r)["id"])

	if oid == nil {
		return http.StatusBadRequest, nil
	}

	err := s.materialRepository.DeleteMaterial(oid)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, oid
}

func (s *service) AddMaterialToBudget(r *http.Request) int {
	budgetId := core.ConvertHexToObjectId(mux.Vars(r)["budgetId"])
	materialId := core.ConvertHexToObjectId(mux.Vars(r)["materialId"])

	if budgetId == nil || materialId == nil {
		return http.StatusBadRequest
	}

	budgetDTO := s.budgetRepository.FindBudgetByOID(budgetId)

	if budgetDTO == nil {
		return http.StatusNotFound
	}

	materialDTO := s.materialRepository.FindMaterialByOID(materialId)

	if materialDTO == nil {
		return http.StatusNotFound
	}

	var materialDetails *MaterialDetailsDTO = &MaterialDetailsDTO{}

	invalidBody := core.DecodeBody(r, materialDetails)

	if invalidBody {
		return http.StatusBadRequest
	}

	err := validate(materialDTO, materialDetails)

	if err != nil {
		return http.StatusBadRequest
	}

	dimension := getMaterialDimension(materialDetails.Metric, materialDTO.Dimensions)

	var budgetMaterial *budget.BudgetMaterial = &budget.BudgetMaterial{
		ID:       primitive.NewObjectID(),
		Quantity: materialDetails.Quantity,
		Name:     materialDTO.Name,
		Dimension: budget.BudgetMaterialDimension{
			ID:       dimension.ID,
			Metric:   dimension.Metric,
			Quantity: dimension.Quantity,
			Price:    dimension.Price,
		},
		Price: float64(materialDetails.Quantity) / dimension.Quantity * dimension.Price,
	}

	err = s.budgetRepository.AddMaterialToBudget(budgetId, budgetMaterial)

	if err != nil {
		return http.StatusInternalServerError
	}

	err = s.budgetRepository.UpdateBudgetByIdPrice(budgetId)

	if err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (s *service) ChangeMaterialPrice(r *http.Request) (int, *MaterialDTO) {
	materialDimensionId := mux.Vars(r)["id"]
	materialDimensionOid := core.ConvertHexToObjectId(materialDimensionId)

	if materialDimensionOid == nil {
		return http.StatusBadRequest, nil
	}

	var materialDimensionPrice *MaterialDimensionPriceDTO = &MaterialDimensionPriceDTO{}

	invalidBody := core.DecodeBody(r, materialDimensionPrice)

	if invalidBody {
		return http.StatusBadRequest, nil
	}

	err := s.materialRepository.ChangeMaterialPrice(materialDimensionOid, materialDimensionPrice)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	materialUpdated := s.materialRepository.FindMaterialByPackageId(materialDimensionOid)

	if materialUpdated == nil {
		return http.StatusInternalServerError, nil
	}

	budget := s.budgetRepository.FindBudgetsByDimensionId(materialDimensionOid)

	if len(budget) == 0 {
		return http.StatusOK, materialUpdated
	}

	err = s.budgetRepository.UpdateMaterialDimensionPrice(materialDimensionOid, materialDimensionPrice.Price)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	for i := 0; i < len(budget); i++ {
		var recipePrice float64 = 0
		for j := 0; j < len(budget[i].Materials); j++ {
			budget[i].Materials[j].Price = budget[i].Materials[j].Quantity / budget[i].Materials[j].Dimension.Quantity * budget[i].Materials[j].Dimension.Price
			recipePrice += budget[i].Materials[j].Price
		}
		budget[i].Price = recipePrice * 3

		err := s.budgetRepository.UpdateMaterialsPrice(materialDimensionOid, budget[i])

		if err != nil {
			log.Println(err.Error())
		}
	}

	return http.StatusOK, materialUpdated
}

func validate(Material *MaterialDTO, MaterialDetails *MaterialDetailsDTO) error {
	if !MaterialMetricMatches(MaterialDetails.Metric, Material.Dimensions) {
		log.Println("La unidad de medida no coincide")
		return errors.New("la unidad de medida no coincide")
	}

	if MaterialDetails.Quantity == 0 {
		log.Println("La cantidad del Materiale no puede ser 0")
		return errors.New("la cantidad del Materiale no puede ser 0")
	}
	return nil
}

func MaterialMetricMatches(metric string, dimensions []DimensionDTO) bool {
	for _, dimension := range dimensions {
		if fmt.Sprintf("%g %s", dimension.Quantity, dimension.Metric) == metric {
			return true
		}
	}
	return false
}

func getMaterialDimension(metric string, dimensions []DimensionDTO) *DimensionDTO {
	for _, dimension := range dimensions {
		if fmt.Sprintf("%g %s", dimension.Quantity, dimension.Metric) == metric {
			return &dimension
		}
	}
	return nil
}
