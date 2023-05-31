// Package routes contains all routes of the application
package routes

import (
	menuController "github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers/menu"
	"github.com/gin-gonic/gin"
)

// MenuRoutes is a function that contains all menu routes
func MenuRoutes(router *gin.RouterGroup, controller *menuController.Controller) {

	routerMenu := router.Group("/menu")
	{
		routerMenu.POST("/", controller.NewMenu)
		routerMenu.GET("/:id", controller.GetMenusByID)
		routerMenu.GET("/", controller.GetAllMenus)
		routerMenu.DELETE("/:id", controller.DeleteMenu)
	}

}
