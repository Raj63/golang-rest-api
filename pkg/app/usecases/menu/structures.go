// Package menu provides the use case for menu
package menu

import (
	domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"
)

// NewMenu is a struct that contains the data for a new menu
type NewMenu struct {
	Name        string `json:"name" example:"Paracetamol"`
	Description string `json:"description" example:"Some Description"`
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
