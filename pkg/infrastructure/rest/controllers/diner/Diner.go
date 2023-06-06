// Package diner contains the diner controller
package diner

import (
	"errors"

	useCaseDiner "github.com/Raj63/golang-rest-api/pkg/app/usecases/diner"
	domainDiner "github.com/Raj63/golang-rest-api/pkg/domain/diner"
	domainErrors "github.com/Raj63/golang-rest-api/pkg/domain/errors"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller is a struct that contains the diner service
type Controller struct {
	DinerService useCaseDiner.Service
}

// NewDiner godoc
//
//	@Tags			diners
//	@Summary		Create New Diner
//	@Description	Create new diner on the system
//	@Accept			json
//	@Produce		json
//	@Param			data	body		NewDinerRequest	true	"body data"
//	@Success		201		{object}	domainDiner.Diner
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/diners [post]
func (c *Controller) NewDiner(ctx *gin.Context) {
	var request NewDinerRequest

	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	newDiner := useCaseDiner.NewDiner{
		Name:        request.Name,
		TableNumber: request.TableNumber,
	}

	var result *domainDiner.Diner
	var err error

	result, err = c.DinerService.Create(ctx.Request.Context(), &newDiner)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// GetAllDiners godoc
//
//	@Tags			diners
//	@Summary		Get all Diners
//	@Description	Get all Diners on the system
//	@Param			limit	query		int64	true	"limit"
//	@Param			page	query		int64	true	"page"
//	@Success		200		{object}	[]useCaseDiner.PaginationResultDiner
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/diners [get]
func (c *Controller) GetAllDiners(ctx *gin.Context) {
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

	diners, err := c.DinerService.GetAll(ctx.Request.Context(), page, limit)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}
	ctx.JSON(http.StatusOK, diners)
}

// GetDinersByID godoc
//
//	@Tags			diners
//	@Summary		Get diners by ID
//	@Description	Get Diners by ID on the system
//	@Param			diner_id	path		int64	true	"id of diner"
//	@Success		200			{object}	domainDiner.Diner
//	@Failure		400			{object}	MessageResponse
//	@Failure		500			{object}	MessageResponse
//	@Router			/diners/{diner_id} [get]
func (c *Controller) GetDinersByID(ctx *gin.Context) {
	dinerID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("diner id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainDiner, err := c.DinerService.GetByID(ctx.Request.Context(), dinerID)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}

	ctx.JSON(http.StatusOK, domainDiner)
}

// DeleteDiner is the controller to delete a diner
//
//	@Tags			diners
//	@Summary		Delete diners by ID
//	@Description	Delete Diners by ID on the system
//	@Param			diner_id	path		int64	true	"id of diner"
//	@Success		200			{object}	MessageResponse
//	@Failure		400			{object}	MessageResponse
//	@Failure		500			{object}	MessageResponse
//	@Router			/diners/{diner_id} [delete]
func (c *Controller) DeleteDiner(ctx *gin.Context) {
	dinerID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param diner id is necessary in the url"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	err = c.DinerService.Delete(ctx.Request.Context(), dinerID)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}
