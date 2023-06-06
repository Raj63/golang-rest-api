// Package order provides the use case for orders
package order

// NewOrder is a struct that contains the data for a new Order
type NewOrder struct {
	DinnerID int64 `db:"diner_id" example:"1"`
	MenuID   int64 `db:"menu_id" example:"3"`
	Quantity int   `db:"quantity" example:"2"`
}
