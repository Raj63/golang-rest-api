// Package adapter is a layer that connects the infrastructure with the application layer
package adapter

import (
	menuService "github.com/Raj63/golang-rest-api/pkg/app/usecases/menu"
	menuRepository "github.com/Raj63/golang-rest-api/pkg/infrastructure/repository/menu"
	menuController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/menu"
	sdksql "github.com/Raj63/golang-rest-api/pkg/infrastructure/sql"
)

// MenuAdapter is a function that returns a menu controller
func MenuAdapter(db *sdksql.DB) *menuController.Controller {
	mRepository := menuRepository.Repository{DB: db}
	service := menuService.Service{MenuRepository: mRepository}
	return &menuController.Controller{MenuService: service}
}
