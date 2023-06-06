// Package diner contains the repository implementation for the diner entity
package diner

import (
	"time"
)

// Diner is a struct that contains the diner model
type Diner struct {
	ID          int64     `db:"id" example:"123"`
	TableNumber int       `db:"table_no" example:"101"`
	Name        string    `db:"name" example:"Mr. Smith"`
	CreatedAt   time.Time `db:"created_at" example:"2021-02-24 20:19:39"`
	UpdatedAt   time.Time `db:"updated_at" example:"2021-02-24 20:19:39"`
}

// TableName overrides the table name used by User to `users`
func (*Diner) TableName() string {
	return "diners"
}
