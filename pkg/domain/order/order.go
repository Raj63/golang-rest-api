// Package order contains the business logic for the order entity
package order

import (
	"time"
)

// Request is a struct that contains the request order information
type Request struct {
	ID        int64     `json:"id" example:"123"`
	DinnerID  int64     `json:"diner_id" example:"1"`
	MenuID    int64     `json:"menu_id" example:"3"`
	Quantity  int       `json:"quantity" example:"2"`
	CreatedAt time.Time `json:"created_at,omitempty" `
}

// Response is a struct that contains the response order information
type Response struct {
	ID              int64     `json:"id" example:"123"`
	DinnerName      string    `json:"diner_name" example:"Mr. Smith"`
	MenuName        string    `json:"menu_name" example:"HCDB"`
	MenuDescription string    `json:"menu_description" example:"Hyderabadi Chicken Dum Briyani"`
	Quantity        int       `json:"quantity" example:"2"`
	CreatedAt       time.Time `json:"created_at,omitempty" `
	UpdatedAt       time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39"`
}

// Service is a interface that contains the methods for the order service
type Service interface {
	Get(int) ([]*Response, error)
	Create(*Request) error
	Delete(int) error
}
