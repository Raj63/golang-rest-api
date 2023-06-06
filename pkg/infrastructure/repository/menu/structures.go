// Package menu contains the repository implementation for the menu entity
package menu

import (
	"time"
)

// Menu is a struct that contains the menu model
type Menu struct {
	ID          int64     `db:"id" example:"123"`
	Name        string    `db:"name" example:"Hyderabadi Dum Briyani"`
	Description string    `db:"description" example:"Some Description"`
	Price       int       `db:"price" example:"20050"`
	CreatedAt   time.Time `db:"created_at" example:"2021-02-24 20:19:39"`
	UpdatedAt   time.Time `db:"updated_at" example:"2021-02-24 20:19:39"`
	Count       int       `db:"count" example:"3"`
}

// TableName overrides the table name used by User to `users`
func (*Menu) TableName() string {
	return "menus"
}
