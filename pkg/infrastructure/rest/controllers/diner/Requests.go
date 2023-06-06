// Package diner contains the diner controller
package diner

// NewDinerRequest is a struct that contains the new diner request information
type NewDinerRequest struct {
	Name        string `json:"name" example:"Paracetamol" binding:"required"`
	TableNumber int    `json:"table_no" example:"101"`
}
