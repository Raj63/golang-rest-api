// Package order contains the order controller
package order

import (
	"errors"

	useCaseOrder "github.com/Raj63/golang-rest-api/pkg/app/usecases/order"
	domainErrors "github.com/Raj63/golang-rest-api/pkg/domain/errors"
	domainOrder "github.com/Raj63/golang-rest-api/pkg/domain/order"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/rest/controllers"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Controller is a struct that contains the order service
type Controller struct {
	OrderService useCaseOrder.Service
}

// NewOrder godoc
//
//	@Tags			orders
//	@Summary		Create New order
//	@Description	Create new order on the system
//	@Accept			json
//	@Produce		json
//	@Param			data	body		NewOrderRequest	true	"body data"
//	@Success		201		{object}	domainOrder.Request
//	@Failure		400		{object}	MessageResponse
//	@Failure		500		{object}	MessageResponse
//	@Router			/orders [post]
func (c *Controller) NewOrder(ctx *gin.Context) {
	var request NewOrderRequest

	if err := controllers.BindJSON(ctx, &request); err != nil {
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	newOrder := useCaseOrder.NewOrder{
		DinnerID: request.DinnerID,
		MenuID:   request.MenuID,
		Quantity: request.Quantity,
	}

	var result *domainOrder.Request
	var err error

	result, err = c.OrderService.Create(ctx.Request.Context(), &newOrder)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// GetOrdersByDinerID godoc
//
//	@Tags			orders
//	@Summary		Get orders by Diner ID
//	@Description	Get orders by Diner ID on the system
//	@Param			diner_id	path		int64	true	"id of diner"
//	@Success		200			{object}	[]domainOrder.Response
//	@Failure		400			{object}	MessageResponse
//	@Failure		500			{object}	MessageResponse
//	@Router			/orders/{diner_id} [get]
func (c *Controller) GetOrdersByDinerID(ctx *gin.Context) {
	orderID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("diner id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	domainOrders, err := c.OrderService.GetByID(ctx.Request.Context(), orderID)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}

	ctx.JSON(http.StatusOK, domainOrders)
}

// DeleteOrder is the controller to delete a order
//
//	@Tags			orders
//	@Summary		Delete orders by ID
//	@Description	Delete orders by ID on the system
//	@Param			order_id	path		int64	true	"id of order"
//	@Success		200			{object}	MessageResponse
//	@Failure		400			{object}	MessageResponse
//	@Failure		500			{object}	MessageResponse
//	@Router			/orders/{order_id} [delete]
func (c *Controller) DeleteOrder(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		appError := domainErrors.NewAppError(errors.New("param order id is necessary in the url"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	err = c.OrderService.Delete(ctx.Request.Context(), orderID)
	if err != nil {
		_ = ctx.Error(domainErrors.NewAppError(err, domainErrors.RepositoryError))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}
