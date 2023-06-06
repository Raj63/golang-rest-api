// Package diner provides the use case for diner
package diner

import (
	domainDiner "github.com/Raj63/golang-rest-api/pkg/domain/diner"
)

func (n *NewDiner) toDomainMapper() *domainDiner.Diner {
	return &domainDiner.Diner{
		Name:        n.Name,
		TableNumber: n.TableNumber,
	}
}
