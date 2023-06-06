// Package routes contains all routes of the application
package routes

import (
	orderController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/order"
	"github.com/gin-gonic/gin"
)

// OrderRoutes is a function that contains all order routes
func OrderRoutes(router *gin.RouterGroup, controller *orderController.Controller) {

	routerOrder := router.Group("/orders")
	{
		routerOrder.POST("/", controller.NewOrder)
		routerOrder.GET("/:id", controller.GetOrdersByDinerID)
		routerOrder.DELETE("/:id", controller.DeleteOrder)
	}

}
