// Package routes contains all routes of the application
package routes

import (
	dinerController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/diner"
	"github.com/gin-gonic/gin"
)

// DinerRoutes is a function that contains all diner routes
func DinerRoutes(router *gin.RouterGroup, controller *dinerController.Controller) {

	routerDiner := router.Group("/diners")
	{
		routerDiner.POST("/", controller.NewDiner)
		routerDiner.GET("/:id", controller.GetDinersByID)
		routerDiner.GET("/", controller.GetAllDiners)
		routerDiner.DELETE("/:id", controller.DeleteDiner)
	}

}
