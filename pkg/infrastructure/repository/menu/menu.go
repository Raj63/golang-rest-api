// Package menu contains the repository implementation for the menu entity
package menu

import (
	domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"
	sdksql "github.com/Raj63/golang-rest-api/pkg/infrastructure/sql"
)

// Repository is a struct that contains the database implementation for menu entity
type Repository struct {
	DB *sdksql.DB
}

// GetAll Fetch all menu data
func (r *Repository) GetAll(page int64, limit int64) (*PaginationResultMenu, error) {
	var menus []Menu
	var total int64

	offset := (page - 1) * limit
	_ = offset

	// Read ALL menu from DB
	numPages := (total + limit - 1) / limit
	var nextCursor, prevCursor uint
	if page < numPages {
		nextCursor = uint(page + 1)
	}
	if page > 1 {
		prevCursor = uint(page - 1)
	}

	return &PaginationResultMenu{
		Data:       arrayToDomainMapper(&menus),
		Total:      total,
		Limit:      limit,
		Current:    page,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		NumPages:   numPages,
	}, nil
}

// Create ... Insert New data
func (r *Repository) Create(newMenu *domainMenu.Menu) (*domainMenu.Menu, error) {
	menu := fromDomainMapper(newMenu)
	// store into DB
	createdMenu := menu.toDomainMapper()
	return createdMenu, nil
}

// GetByID ... Fetch only one menu by Id
func (r *Repository) GetByID(id int) (*domainMenu.Menu, error) {
	var menu Menu

	return menu.toDomainMapper(), nil
}

// Delete ... Delete menu
func (r *Repository) Delete(id int) (err error) {
	// tx := r.DB.Delete(&domainMenu.Menu{}, id)
	// if tx.Error != nil {
	// 	err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	// 	return
	// }

	// if tx.RowsAffected == 0 {
	// 	err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	// }

	return
}
