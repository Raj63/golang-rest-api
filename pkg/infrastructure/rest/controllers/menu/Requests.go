// Package menu contains the menu controller
package menu

// NewMenuRequest is a struct that contains the new menu request information
type NewMenuRequest struct {
	Name        string  `json:"name" example:"Paracetamol" binding:"required"`
	Description string  `json:"description" example:"Something" binding:"required"`
	Price       float64 `json:"price" example:"200.50" binding:"required"`
}
