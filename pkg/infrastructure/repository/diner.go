package repository

import (
	"context"

	domainDiner "github.com/Raj63/golang-rest-api/pkg/domain/diner"
)

// Diners specifies the repository contracts
type Diners interface {
	GetTotalCount(ctx context.Context) (int64, error)
	GetAll(ctx context.Context, page int64, limit int64) (*PaginationResultDiner, error)
	Create(ctx context.Context, newDiner *domainDiner.Diner) (*domainDiner.Diner, error)
	GetByID(ctx context.Context, id int64) (*domainDiner.Diner, error)
	Delete(ctx context.Context, id int64) (err error)
}

// PaginationResultDiner is a struct that contains the pagination result for diner
type PaginationResultDiner struct {
	Data       *[]domainDiner.Diner
	Total      int64
	Limit      int64
	Current    int64
	NextCursor uint
	PrevCursor uint
	NumPages   int64
}
