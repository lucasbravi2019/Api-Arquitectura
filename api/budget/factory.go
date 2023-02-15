package budget

import (
	"github.com/lucasbravi2019/arquitectura/core"
)

func GetBudgetHandlerInstance() *handler {
	if budgetHandlerInstance == nil {
		budgetHandlerInstance = &handler{
			service: GetBudgetServiceInstance(),
		}
	}
	return budgetHandlerInstance
}

func GetBudgetServiceInstance() *service {
	if budgetServiceInstance == nil {
		budgetServiceInstance = &service{
			budgetRepository: GetBudgetRepositoryInstance(),
		}
	}
	return budgetServiceInstance
}

func GetBudgetRepositoryInstance() *repository {
	if budgetRepositoryInstance == nil {
		budgetRepositoryInstance = &repository{
			db: core.GetDatabaseConnection().Collection("budgets"),
		}
	}
	return budgetRepositoryInstance
}
