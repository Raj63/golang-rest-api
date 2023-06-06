// Package diner provides the use case for diner
package diner

import (
	domainDiner "github.com/Raj63/golang-rest-api/pkg/domain/diner"
)

// NewDiner is a struct that contains the data for a new diner
type NewDiner struct {
	Name        string `json:"name" example:"Mr. Smith"`
	TableNumber int    `json:"table_no" example:"101"`
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
