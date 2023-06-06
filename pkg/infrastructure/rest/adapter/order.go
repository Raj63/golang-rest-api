// Package adapter is a layer that connects the infrastructure with the application layer
package adapter

import (
	orderService "github.com/Raj63/golang-rest-api/pkg/app/usecases/order"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	orderRepository "github.com/Raj63/golang-rest-api/pkg/infrastructure/repository/order"
	orderController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/order"
	sdksql "github.com/Raj63/golang-rest-api/pkg/infrastructure/sql"
)

// OrderAdapter is a function that returns a order controller
func OrderAdapter(db *sdksql.DB, logger *logger.Logger) *orderController.Controller {
	mRepository := orderRepository.Repository{Store: db, Logger: logger}
	service := orderService.Service{OrderRepository: mRepository}
	return &orderController.Controller{OrderService: service}
}
