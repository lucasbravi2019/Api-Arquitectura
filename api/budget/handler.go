package budget

import (
	"net/http"

	"github.com/lucasbravi2019/arquitectura/core"
)

type handler struct {
	service BudgetService
}

var budgetHandlerInstance *handler

type RecipeHandler interface {
	GetAllBudgets(w http.ResponseWriter, r *http.Request)
	GetBudget(w http.ResponseWriter, r *http.Request)
	CreateBudget(w http.ResponseWriter, r *http.Request)
	UpdateBudgetName(w http.ResponseWriter, r *http.Request)
	DeleteBudget(w http.ResponseWriter, r *http.Request)

	GetBudgetRoutes() core.Routes
}

func (h *handler) GetAllBudgets(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.GetAllBudgets()
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) GetBudget(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.GetBudget(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.CreateBudget(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) UpdateBudgetName(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.UpdateBudgetName(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.DeleteBudget(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) GetBudgetRoutes() core.Routes {
	return core.Routes{
		core.Route{
			Path:        "/budgets",
			HandlerFunc: h.GetAllBudgets,
			Method:      "GET",
		},
		core.Route{
			Path:        "/budgets",
			HandlerFunc: h.CreateBudget,
			Method:      "POST",
		},
		core.Route{
			Path:        "/budgets/{id}",
			HandlerFunc: h.UpdateBudgetName,
			Method:      "PUT",
		},
		core.Route{
			Path:        "/budgets/{id}",
			HandlerFunc: h.GetBudget,
			Method:      "GET",
		},
		core.Route{
			Path:        "/budgets/{id}",
			HandlerFunc: h.DeleteBudget,
			Method:      "DELETE",
		},
	}
}
