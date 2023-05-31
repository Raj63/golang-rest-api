// Package menu contains the menu controller
package menu

import (
	"errors"

	useCaseMenu "github.com/Raj63/golang-rest-api/pkg/app/usecases/menu"
	domainErrors "github.com/Raj63/golang-rest-api/pkg/domain/errors"
	domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller is a struct that contains the menu service
type Controller struct {
	MenuService useCaseMenu.Service
}

// NewMenu godoc
//
//	@Tags			menu
//	@Summary		Create New Menu
//	@Description	Create new menu on the system
//	@Accept			json
//	@Produce		json
//	@Param			data	body		NewMenuRequest	true	"body data"
//	@Success		200		{object}	domainMenu.Menu
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/menu [post]
func (c *Controller) NewMenu(ctx *gin.Context) {
	var request NewMenuRequest

	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	newMenu := useCaseMenu.NewMenu{
		Name:        request.Name,
		Description: request.Description,
	}

	var result *domainMenu.Menu
	var err error

	result, err = c.MenuService.Create(&newMenu)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetAllMenus godoc
//
//	@Tags			menu
//	@Summary		Get all Menus
//	@Description	Get all Menus on the system
//	@Param			limit	query		string	true	"limit"
//	@Param			page	query		string	true	"page"
//	@Success		200		{object}	[]useCaseMenu.PaginationResultMenu
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/menu [get]
func (c *Controller) GetAllMenus(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "20")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param page is necessary to be an integer"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param limit is necessary to be an integer"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	menus, err := c.MenuService.GetAll(page, limit)
	if err != nil {
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	ctx.JSON(http.StatusOK, menus)
}

// GetMenusByID godoc
//
//	@Tags			menu
//	@Summary		Get menus by ID
//	@Description	Get Menus by ID on the system
//	@Param			menu_id	path		int	true	"id of menu"
//	@Success		200			{object}	domainMenu.Menu
//	@Failure		400			{object}	MessageResponse
//	@Failure		500			{object}	MessageResponse
//	@Router			/menu/{menu_id} [get]
func (c *Controller) GetMenusByID(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("menu id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainMenu, err := c.MenuService.GetByID(menuID)
	if err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	ctx.JSON(http.StatusOK, domainMenu)
}

// DeleteMenu is the controller to delete a menu
func (c *Controller) DeleteMenu(ctx *gin.Context) {
	menuID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param id is necessary in the url"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	err = c.MenuService.Delete(menuID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}
