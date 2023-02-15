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
	MaterialOid := core.ConvertHexToObjectId(mux.Vars(r)["materialId"])

	if budgetId == nil || MaterialOid == nil {
		return http.StatusBadRequest
	}

	recipe := s.budgetRepository.FindBudgetByOID(budgetId)

	if recipe == nil {
		return http.StatusNotFound
	}

	MaterialDTO := s.materialRepository.FindMaterialByOID(MaterialOid)

	if MaterialDTO == nil {
		return http.StatusNotFound
	}

	var MaterialDetails *MaterialDetailsDTO = &MaterialDetailsDTO{}

	invalidBody := core.DecodeBody(r, MaterialDetails)

	if invalidBody {
		return http.StatusBadRequest
	}

	err := validate(MaterialDTO, MaterialDetails)

	if err != nil {
		return http.StatusBadRequest
	}

	envase := getMaterialPackage(MaterialDetails.Metric, MaterialDTO.Dimensions)

	var recipeMaterial *budget.BudgetMaterial = &budget.BudgetMaterial{
		ID:       primitive.NewObjectID(),
		Quantity: MaterialDetails.Quantity,
		Name:     MaterialDTO.Name,
		Dimension: budget.BudgetMaterialDimension{
			ID:       envase.ID,
			Metric:   envase.Metric,
			Quantity: envase.Quantity,
			Price:    envase.Price,
		},
		Price: float64(MaterialDetails.Quantity) / envase.Quantity * envase.Price,
	}

	err = s.budgetRepository.AddMaterialToBudget(budgetId, recipeMaterial)

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
	MaterialPackageId := mux.Vars(r)["id"]
	MaterialPackageOid := core.ConvertHexToObjectId(MaterialPackageId)

	if MaterialPackageOid == nil {
		return http.StatusBadRequest, nil
	}

	var MaterialPackagePrice *MaterialDimensionPriceDTO = &MaterialDimensionPriceDTO{}

	invalidBody := core.DecodeBody(r, MaterialPackagePrice)

	if invalidBody {
		return http.StatusBadRequest, nil
	}

	err := s.materialRepository.ChangeMaterialPrice(MaterialPackageOid, MaterialPackagePrice)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	MaterialUpdated := s.materialRepository.FindMaterialByPackageId(MaterialPackageOid)

	if MaterialUpdated == nil {
		return http.StatusInternalServerError, nil
	}

	budget := s.budgetRepository.FindBudgetsByDimensionId(MaterialPackageOid)

	if len(budget) == 0 {
		return http.StatusOK, MaterialUpdated
	}

	err = s.budgetRepository.UpdateMaterialDimensionPrice(MaterialPackageOid, MaterialPackagePrice.Price)

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

		err := s.budgetRepository.UpdateMaterialsPrice(MaterialPackageOid, budget[i])

		if err != nil {
			log.Println(err.Error())
		}
	}

	return http.StatusOK, MaterialUpdated
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

func MaterialMetricMatches(metric string, packages []DimensionDTO) bool {
	for _, pack := range packages {
		if fmt.Sprintf("%g %s", pack.Quantity, pack.Metric) == metric {
			return true
		}
	}
	return false
}

func getMaterialPackage(metric string, packages []DimensionDTO) *DimensionDTO {
	for _, pack := range packages {
		if fmt.Sprintf("%g %s", pack.Quantity, pack.Metric) == metric {
			return &pack
		}
	}
	return nil
}
