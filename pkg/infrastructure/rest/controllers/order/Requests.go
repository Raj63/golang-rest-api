// Package order contains the order controller
package order

// NewOrderRequest is a struct that contains the new order request information
type NewOrderRequest struct {
	DinnerID int64 `json:"diner_id" example:"1" binding:"required"`
	MenuID   int64 `json:"menu_id" example:"3" binding:"required"`
	Quantity int   `json:"quantity" example:"2" binding:"required"`
}
