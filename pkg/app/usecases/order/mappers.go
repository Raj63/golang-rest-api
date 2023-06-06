// Package order provides the use case for order
package order

import (
	domainOrder "github.com/Raj63/golang-rest-api/pkg/domain/order"
)

func (n *NewOrder) toDomainMapper() *domainOrder.Request {
	return &domainOrder.Request{
		DinnerID: n.DinnerID,
		MenuID:   n.MenuID,
		Quantity: n.Quantity,
	}
}
