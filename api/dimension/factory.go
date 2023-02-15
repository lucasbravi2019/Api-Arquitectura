package dimension

import (
	"github.com/lucasbravi2019/arquitectura/api/budget"
	"github.com/lucasbravi2019/arquitectura/api/material"
	"github.com/lucasbravi2019/arquitectura/core"
)

func GetDimensionHandlerInstance() *handler {
	if dimensionHandlerInstance == nil {
		dimensionHandlerInstance = &handler{
			service: GetDimensionServiceInstance(),
		}
	}
	return dimensionHandlerInstance
}

func GetDimensionServiceInstance() *service {
	if dimensionServiceInstance == nil {
		dimensionServiceInstance = &service{
			dimensionRepository: GetDimensionRepositoryInstance(),
			materialRepository:  material.GetMaterialRepositoryInstance(),
			budgetRepository:    budget.GetBudgetRepositoryInstance(),
		}
	}
	return dimensionServiceInstance
}

func GetDimensionRepositoryInstance() *repository {
	if dimensionRepositoryInstance == nil {
		dimensionRepositoryInstance = &repository{
			db: core.GetDatabaseConnection().Collection("dimensions"),
		}
	}
	return dimensionRepositoryInstance
}
