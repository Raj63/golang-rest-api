// Package adapter is a layer that connects the infrastructure with the application layer
package adapter

import (
	menuService "github.com/Raj63/golang-rest-api/pkg/app/usecases/menu"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	menuRepository "github.com/Raj63/golang-rest-api/pkg/infrastructure/repository/menu"
	menuController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/menu"
	sdksql "github.com/Raj63/golang-rest-api/pkg/infrastructure/sql"
)

// MenuAdapter is a function that returns a menu controller
func MenuAdapter(db *sdksql.DB, logger *logger.Logger) *menuController.Controller {
	mRepository := menuRepository.Repository{Store: db, Logger: logger}
	service := menuService.Service{MenuRepository: mRepository}
	return &menuController.Controller{MenuService: service}
}
