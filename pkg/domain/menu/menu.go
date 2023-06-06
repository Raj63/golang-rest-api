// Package menu contains the business logic for the menu entity
package menu

import (
	"time"
)

// Menu is a struct that contains the menu information
type Menu struct {
	ID          int64     `json:"id" example:"123"`
	Name        string    `json:"name" example:"Hyderabadi Dum Briyani"`
	Description string    `json:"description" example:"Some Description"`
	Price       float64   `json:"price" example:"200.50"`
	CreatedAt   time.Time `json:"created_at,omitempty" `
	UpdatedAt   time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39"`
}

// Service is a interface that contains the methods for the menu service
type Service interface {
	Get(int) (*Menu, error)
	GetAll() ([]*Menu, error)
	Create(*Menu) error
	Delete(int) error
}
