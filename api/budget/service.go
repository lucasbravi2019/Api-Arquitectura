package budget

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lucasbravi2019/arquitectura/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	budgetRepository BudgetRepository
}

type BudgetService interface {
	GetAllBudgets() (int, *[]BudgetDTO)
	GetBudget(r *http.Request) (int, *BudgetDTO)
	CreateBudget(r *http.Request) (int, *BudgetDTO)
	UpdateBudgetName(r *http.Request) (int, *BudgetDTO)
	DeleteBudget(r *http.Request) (int, *primitive.ObjectID)
}

var budgetServiceInstance *service

func (s *service) GetAllBudgets() (int, *[]BudgetDTO) {
	recipes := s.budgetRepository.FindAllBudgets()
	return http.StatusOK, recipes
}

func (s *service) GetBudget(r *http.Request) (int, *BudgetDTO) {
	oid := core.ConvertHexToObjectId(mux.Vars(r)["id"])

	if oid == nil {
		return http.StatusBadRequest, nil
	}

	budget := s.budgetRepository.FindBudgetByOID(oid)

	if budget == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, budget
}

func (s *service) CreateBudget(r *http.Request) (int, *BudgetDTO) {
	var budgetName *BudgetNameDTO = &BudgetNameDTO{}

	invalidBody := core.DecodeBody(r, budgetName)

	if invalidBody {
		return http.StatusBadRequest, nil
	}

	oid := s.budgetRepository.CreateBudget(budgetName)

	if oid == nil {
		return http.StatusInternalServerError, nil
	}

	budget := s.budgetRepository.FindBudgetByOID(oid)

	if budget == nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusCreated, budget
}

func (s *service) UpdateBudgetName(r *http.Request) (int, *BudgetDTO) {
	oid := core.ConvertHexToObjectId(mux.Vars(r)["id"])

	if oid == nil {
		return http.StatusBadRequest, nil
	}

	var budget *BudgetNameDTO = &BudgetNameDTO{}

	invalidBody := core.DecodeBody(r, budget)

	if invalidBody {
		return http.StatusBadRequest, nil
	}

	err := s.budgetRepository.UpdateBudgetName(oid, budget)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	budgetUpdated := s.budgetRepository.FindBudgetByOID(oid)

	if budgetUpdated == nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, budgetUpdated
}

func (s *service) DeleteBudget(r *http.Request) (int, *primitive.ObjectID) {
	oid := core.ConvertHexToObjectId(mux.Vars(r)["id"])

	if oid == nil {
		return http.StatusBadRequest, nil
	}

	err := s.budgetRepository.DeleteBudget(oid)

	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, oid
}
