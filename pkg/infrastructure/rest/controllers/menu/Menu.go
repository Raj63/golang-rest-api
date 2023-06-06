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
//	@Tags			menus
//	@Summary		Create New Menu
//	@Description	Create new menu on the system
//	@Accept			json
//	@Produce		json
//	@Param			data	body		NewMenuRequest	true	"body data"
//	@Success		201		{object}	domainMenu.Menu
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/menus [post]
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
		Price:       request.Price,
	}

	var result *domainMenu.Menu
	var err error

	result, err = c.MenuService.Create(ctx.Request.Context(), &newMenu)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// GetAllMenus godoc
//
//	@Tags			menus
//	@Summary		Get all Menus
//	@Description	Get all Menus on the system
//	@Param			limit	query		int64	true	"limit"
//	@Param			page	query		int64	true	"page"
//	@Success		200		{object}	[]useCaseMenu.PaginationResultMenu
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/menus [get]
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

	menus, err := c.MenuService.GetAll(ctx.Request.Context(), page, limit)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}
	ctx.JSON(http.StatusOK, menus)
}

// GetTopMenus godoc
//
//		@Tags			menus
//		@Summary		Get top menus by count
//		@Description	Get Top Menus by count on the system
//	    @Param          count       query int true     "top count"
//		@Success		200			{object}	[]domainMenu.Menu
//		@Failure		400			{object}	MessageResponse
//		@Failure		500			{object}	MessageResponse
//		@Router			/menus/top [get]
func (c *Controller) GetTopMenus(ctx *gin.Context) {
	count, err := strconv.Atoi(ctx.Query("count"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("menu id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainMenu, err := c.MenuService.GetByTopCount(ctx.Request.Context(), count)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}

	ctx.JSON(http.StatusOK, domainMenu)
}

// GetMenusByID godoc
//
//	@Tags			menus
//	@Summary		Get menus by ID
//	@Description	Get Menus by ID on the system
//	@Param			menu_id	path		int64	true	"id of menu"
//	@Success		200			{object}	domainMenu.Menu
//	@Failure		400			{object}	MessageResponse
//	@Failure		500			{object}	MessageResponse
//	@Router			/menus/{menu_id} [get]
func (c *Controller) GetMenusByID(ctx *gin.Context) {
	menuID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("menu id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainMenu, err := c.MenuService.GetByID(ctx.Request.Context(), menuID)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}

	ctx.JSON(http.StatusOK, domainMenu)
}

// DeleteMenu is the controller to delete a menu
//
//	@Tags			menus
//	@Summary		Delete menus by ID
//	@Description	Delete Menus by ID on the system
//	@Param			menu_id	path		int64	true	"id of menu"
//	@Success		200			{object}	MessageResponse
//	@Failure		400			{object}	MessageResponse
//	@Failure		500			{object}	MessageResponse
//	@Router			/menus/{menu_id} [delete]
func (c *Controller) DeleteMenu(ctx *gin.Context) {
	menuID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param menu id is necessary in the url"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	err = c.MenuService.Delete(ctx.Request.Context(), menuID)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}
