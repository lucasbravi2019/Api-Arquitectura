package material

import (
	"github.com/lucasbravi2019/arquitectura/api/budget"
	"github.com/lucasbravi2019/arquitectura/core"
)

func GetMaterialHandlerInstance() *handler {
	if materialHandlerInstance == nil {
		materialHandlerInstance = &handler{
			service: GetMaterialServiceInstance(),
		}
	}
	return materialHandlerInstance
}

func GetMaterialServiceInstance() *service {
	if materialServiceInstance == nil {
		materialServiceInstance = &service{
			materialRepository: GetMaterialRepositoryInstance(),
			budgetRepository:   budget.GetBudgetRepositoryInstance(),
		}
	}
	return materialServiceInstance
}

func GetMaterialRepositoryInstance() *repository {
	if materialRepositoryInstance == nil {
		materialRepositoryInstance = &repository{
			materialCollection: core.GetDatabaseConnection().Collection("materials"),
		}
	}
	return materialRepositoryInstance
}
