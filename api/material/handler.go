package material

import (
	"net/http"

	"github.com/lucasbravi2019/arquitectura/core"
)

type handler struct {
	service MaterialService
}

type MaterialHandler interface {
	GetAllMaterials(w http.ResponseWriter, r *http.Request)
	CreateMaterial(w http.ResponseWriter, r *http.Request)
	UpdateMaterial(w http.ResponseWriter, r *http.Request)
	DeleteMaterial(w http.ResponseWriter, r *http.Request)
	AddMaterialToRecipe(w http.ResponseWriter, r *http.Request)
	AddPackageToMaterial(w http.ResponseWriter, r *http.Request)
	RemovePackageFromMaterials(w http.ResponseWriter, r *http.Request)
	ChangeMaterialPrice(w http.ResponseWriter, r *http.Request)
	GetMaterialRoutes() core.Routes
}

var materialHandlerInstance *handler

func (h *handler) GetAllMaterials(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.GetAllMaterials()
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.CreateMaterial(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.UpdateMaterial(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.DeleteMaterial(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) AddMaterialToBudget(w http.ResponseWriter, r *http.Request) {
	statusCode := h.service.AddMaterialToBudget(r)
	core.EncodeJsonResponse(w, statusCode, nil)
}

func (h *handler) ChangeMaterialPrice(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.ChangeMaterialPrice(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) GetMaterialRoutes() core.Routes {
	return core.Routes{
		core.Route{
			Path:        "/materials",
			HandlerFunc: h.GetAllMaterials,
			Method:      "GET",
		},
		core.Route{
			Path:        "/materials",
			HandlerFunc: h.CreateMaterial,
			Method:      "POST",
		},
		core.Route{
			Path:        "/materials/{id}",
			HandlerFunc: h.UpdateMaterial,
			Method:      "PUT",
		},
		core.Route{
			Path:        "/materials/{id}/price",
			HandlerFunc: h.ChangeMaterialPrice,
			Method:      "PUT",
		},
		core.Route{
			Path:        "/materials/{id}",
			HandlerFunc: h.DeleteMaterial,
			Method:      "DELETE",
		},
		core.Route{
			Path:        "/materials/{materialId}/budgets/{budgetId}",
			HandlerFunc: h.AddMaterialToBudget,
			Method:      "PUT",
		},
	}
}
