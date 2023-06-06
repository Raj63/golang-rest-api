// Package diner contains the business logic for the diner entity
package diner

import (
	"context"
	"time"
)

// Diner is a struct that contains the diner information
type Diner struct {
	ID          int64     `json:"id" example:"123"`
	Name        string    `json:"name" example:"Mr. Smith"`
	TableNumber int       `json:"table_no" example:"101"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39"`
}

// Service is a interface that contains the methods for the diner service
type Service interface {
	Get(context.Context, int64) (*Diner, error)
	GetAll(context.Context, int64, int64) ([]*Diner, error)
	Create(context.Context, *Diner) error
	Delete(context.Context, int64) error
}
