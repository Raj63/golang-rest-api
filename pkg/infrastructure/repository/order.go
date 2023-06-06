package repository

import (
	"context"

	domainOrder "github.com/Raj63/golang-rest-api/pkg/domain/order"
)

// Orders specifies the repository contracts
type Orders interface {
	Create(ctx context.Context, newOrder *domainOrder.Request) (*domainOrder.Request, error)
	GetByID(ctx context.Context, dinerID int64) ([]domainOrder.Response, error)
	Delete(ctx context.Context, id int) (err error)
}
