package dimension

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lucasbravi2019/arquitectura/api/budget"
	"github.com/lucasbravi2019/arquitectura/api/material"
	"github.com/lucasbravi2019/arquitectura/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	dimensionRepository DimensionRepository
	materialRepository  material.MaterialRepository
	budgetRepository    budget.BudgetRepository
}

type DimensionService interface {
	GetDimensions() (int, *[]Dimension)
	CreateDimension(r *http.Request) (int, *Dimension)
	UpdateDimension(r *http.Request) (int, *Dimension)
	DeleteDimension(r *http.Request) (int, *primitive.ObjectID)
	AddDimensionToMaterial(r *http.Request) int
	RemoveDimensionFromMaterials(r *http.Request) (int, *primitive.ObjectID)
}

var dimensionServiceInstance *service

func (s *service) GetDimensions() (int, *[]Dimension) {
	dimensions := s.dimensionRepository.GetDimensions()

	if dimensions == nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, dimensions
}

func (s *service) CreateDimension(r *http.Request) (int, *Dimension) {
	var dimensionRequest *Dimension = &Dimension{}

	invalidBody := core.DecodeBody(r, dimensionRequest)

	if invalidBody {
		return http.StatusBadRequest, nil
	}

	id := s.dimensionRepository.CreateDimension(dimensionRequest)

	if id == nil {
		return http.StatusInternalServerError, nil
	}

	dimension := s.dimensionRepository.GetDimensionById(id)

	if dimension == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusCreated, dimension
}

func (s *service) UpdateDimension(r *http.Request) (int, *Dimension) {
	oid := core.ConvertHexToObjectId(mux.Vars(r)["id"])

	if oid == nil {
		return http.StatusBadRequest, nil
	}

	var dimensionRequest *Dimension = &Dimension{}

	invalidBody := core.DecodeBody(r, dimensionRequest)

	if invalidBody {
		return http.StatusBadRequest, nil
	}

	err := s.dimensionRepository.UpdateDimension(oid, dimensionRequest)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	dimension := s.dimensionRepository.GetDimensionById(oid)

	if dimension == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, dimension
}

func (s *service) DeleteDimension(r *http.Request) (int, *primitive.ObjectID) {
	oid := core.ConvertHexToObjectId(mux.Vars(r)["id"])

	if oid == nil {
		return http.StatusBadRequest, nil
	}

	err := s.dimensionRepository.DeleteDimension(oid)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	var materialDimension *material.MaterialDimensionDTO = &material.MaterialDimensionDTO{
		DimensionOid: *oid,
	}

	err = s.materialRepository.RemoveDimensionFromMaterials(*materialDimension)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	err = s.budgetRepository.RemoveMaterialByDimensionId(oid)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	err = s.budgetRepository.UpdateBudgetsPrice()

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, oid
}

func (s *service) AddDimensionToMaterial(r *http.Request) int {
	materialOid := mux.Vars(r)["materialId"]
	dimensionOid := mux.Vars(r)["dimensionId"]
	materialId := core.ConvertHexToObjectId(materialOid)
	dimensionId := core.ConvertHexToObjectId(dimensionOid)

	var priceDTO *material.MaterialDimensionPriceDTO = &material.MaterialDimensionPriceDTO{}

	invalidBody := core.DecodeBody(r, priceDTO)

	if invalidBody {
		return http.StatusBadRequest
	}

	envase := s.dimensionRepository.GetDimensionById(dimensionId)

	if envase == nil {
		return http.StatusNotFound
	}

	var materialDimension *material.MaterialDimension = &material.MaterialDimension{
		ID:       envase.ID,
		Metric:   envase.Metric,
		Quantity: envase.Quantity,
		Price:    priceDTO.Price,
	}

	err := s.materialRepository.AddDimensionToMaterial(materialId, dimensionId, materialDimension)

	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (s *service) RemoveDimensionFromMaterials(r *http.Request) (int, *primitive.ObjectID) {
	dimensionId := core.ConvertHexToObjectId(mux.Vars(r)["dimensionId"])

	var ingredientPackageDto *material.MaterialDimensionDTO = &material.MaterialDimensionDTO{
		DimensionOid: *dimensionId,
	}

	err := s.materialRepository.RemoveDimensionFromMaterials(*ingredientPackageDto)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, dimensionId
}
