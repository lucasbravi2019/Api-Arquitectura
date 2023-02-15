package dimension

import (
	"net/http"

	"github.com/lucasbravi2019/arquitectura/core"
)

type handler struct {
	service DimensionService
}

type DimensionHandler interface {
	GetDimensions(w http.ResponseWriter, r *http.Request)
	CreateDimension(w http.ResponseWriter, r *http.Request)
	UpdateDimension(w http.ResponseWriter, r *http.Request)
	DeleteDimension(w http.ResponseWriter, r *http.Request)
	AddDimensionToMaterial(w http.ResponseWriter, r *http.Request)
	RemoveDimensionFromMaterials(w http.ResponseWriter, r *http.Request)
	GetDimensionRoutes() []core.Route
}

var dimensionHandlerInstance *handler

func (h *handler) GetDimensions(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.GetDimensions()
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) CreateDimension(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.CreateDimension(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) UpdateDimension(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.UpdateDimension(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) DeleteDimension(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.DeleteDimension(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) AddDimensionToMaterial(w http.ResponseWriter, r *http.Request) {
	statusCode := h.service.AddDimensionToMaterial(r)
	core.EncodeJsonResponse(w, statusCode, nil)
}

func (h *handler) RemoveDimensionFromMaterials(w http.ResponseWriter, r *http.Request) {
	statusCode, body := h.service.RemoveDimensionFromMaterials(r)
	core.EncodeJsonResponse(w, statusCode, body)
}

func (h *handler) GetDimensionRoutes() core.Routes {
	return core.Routes{
		core.Route{
			Path:        "/dimensions",
			HandlerFunc: h.GetDimensions,
			Method:      "GET",
		},
		core.Route{
			Path:        "/dimensions",
			HandlerFunc: h.CreateDimension,
			Method:      "POST",
		},
		core.Route{
			Path:        "/dimensions/{id}",
			HandlerFunc: h.UpdateDimension,
			Method:      "PUT",
		},
		core.Route{
			Path:        "/dimensions/{id}",
			HandlerFunc: h.DeleteDimension,
			Method:      "DELETE",
		},
		core.Route{
			Path:        "/dimensions/{dimensionId}/materials/{materialId}",
			HandlerFunc: h.AddDimensionToMaterial,
			Method:      "PUT",
		},
		core.Route{
			Path:        "/dimensions/{id}/materials",
			HandlerFunc: h.RemoveDimensionFromMaterials,
			Method:      "DELETE",
		},
	}
}
