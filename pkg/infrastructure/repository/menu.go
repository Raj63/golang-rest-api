package repository

import (
	"context"

	domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"
)

// Menus specifies the repository contracts
type Menus interface {
	GetTotalCount(ctx context.Context) (int64, error)
	GetAll(ctx context.Context, page int64, limit int64) (*PaginationResultMenu, error)
	Create(ctx context.Context, newMenu *domainMenu.Menu) (*domainMenu.Menu, error)
	GetByID(ctx context.Context, id int64) (*domainMenu.Menu, error)
	GetByTopCount(ctx context.Context, count int) ([]domainMenu.Menu, error)
	Delete(ctx context.Context, id int64) (err error)
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
