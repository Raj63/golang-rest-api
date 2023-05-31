// Package routes contains all routes of the application
package routes

import (
	// swaggerFiles for documentation
	_ "github.com/Raj63/golang-rest-api/docs"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/adapter"
	sdksql "github.com/Raj63/golang-rest-api/pkg/infrastructure/sql"
	"github.com/gin-gonic/gin"
)

//	@title			Golang REST APIs
//	@version		2.0
//	@description	Documentation's Golang REST APIs
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Rajesh Kumar Biswas
//	@contact.url	http://github.com/Raj63
//	@contact.email	biswas.rajesh63@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// ApplicationV1Router is a function that contains all routes of the application
//
//	@host		localhost:8080
//	@BasePath	/v1
func ApplicationV1Router(router *gin.Engine, db *sdksql.DB) {
	routerV1 := router.Group("/v1")

	{
		MenuRoutes(routerV1, adapter.MenuAdapter(db))
	}
}
