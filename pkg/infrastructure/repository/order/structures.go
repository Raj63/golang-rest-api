// Package order contains the repository implementation for the order entity
package order

import (
	"time"
)

// Request is a struct that contains the order model
type Request struct {
	ID        int64     `db:"id" example:"123"`
	DinnerID  int64     `db:"diner_id" example:"1"`
	MenuID    int64     `db:"menu_id" example:"3"`
	Quantity  int       `db:"quantity" example:"2"`
	CreatedAt time.Time `db:"created_at" example:"2021-02-24 20:19:39"`
}

// Response is a struct that contains the response order information
type Response struct {
	ID              int64     `db:"id" example:"123"`
	DinnerName      string    `db:"diner_name" example:"Mr. Smith"`
	MenuName        string    `db:"menu_name" example:"HCDB"`
	MenuDescription string    `db:"menu_description" example:"Hyderabadi Chicken Dum Briyani"`
	Quantity        int       `db:"quantity" example:"2"`
	CreatedAt       time.Time `db:"created_at" example:"2021-02-24 20:19:39"`
	UpdatedAt       time.Time `db:"updated_at" example:"2021-02-24 20:19:39"`
}
