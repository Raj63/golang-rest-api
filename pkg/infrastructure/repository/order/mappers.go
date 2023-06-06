// Package order contains the repository implementation for the order entity
package order

import domainOrder "github.com/Raj63/golang-rest-api/pkg/domain/order"

func (order *Request) toDomainMapper() *domainOrder.Request {
	return &domainOrder.Request{
		ID:        order.ID,
		CreatedAt: order.CreatedAt,
	}
}

func fromDomainMapper(order *domainOrder.Request) *Request {
	return &Request{
		ID:        order.ID,
		DinnerID:  order.DinnerID,
		MenuID:    order.MenuID,
		Quantity:  order.Quantity,
		CreatedAt: order.CreatedAt,
	}
}

func (order *Response) toDomainMapper() *domainOrder.Response {
	return &domainOrder.Response{
		ID:              order.ID,
		DinnerName:      order.DinnerName,
		MenuName:        order.MenuName,
		MenuDescription: order.MenuDescription,
		Quantity:        order.Quantity,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func arrayToDomainMapper(orders *[]Response) []domainOrder.Response {
	ordersDomain := make([]domainOrder.Response, len(*orders))
	for i, order := range *orders {
		ordersDomain[i] = *order.toDomainMapper()
	}

	return ordersDomain
}
