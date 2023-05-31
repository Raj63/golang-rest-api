// Package menu contains the repository implementation for the menu entity
package menu

import (
	"time"

	domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"
)

// Menu is a struct that contains the menu model
type Menu struct {
	ID          int       `json:"id" example:"123" gorm:"primaryKey"`
	Name        string    `json:"name" example:"Paracetamol" gorm:"unique"`
	Description string    `json:"description" example:"Some Description"`
	EANCode     string    `json:"ean_code" example:"9900000124" gorm:"unique"`
	Laboratory  string    `json:"laboratory" example:"Roche"`
	CreatedAt   time.Time `json:"created_at,omitempty" example:"2021-02-24 20:19:39" gorm:"autoCreateTime:mili"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" example:"2021-02-24 20:19:39" gorm:"autoUpdateTime:mili"`
}

// TableName overrides the table name used by User to `users`
func (*Menu) TableName() string {
	return "menus"
}

// PaginationResultMenu is a struct that contains the pagination result for menu
type PaginationResultMenu struct {
	Data       *[]domainMenu.Menu
	Total      int64
	Limit      int64
	Current    int64
	NextCursor uint
	PrevCursor uint
	NumPages   int64
}
