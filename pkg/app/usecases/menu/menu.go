// Package menu provides the use case for menu
package menu

import (
	menuDomain "github.com/Raj63/golang-rest-api/pkg/domain/menu"
	menuRepository "github.com/Raj63/golang-rest-api/pkg/infrastructure/repository/menu"
)

// Service is a struct that contains the repository implementation for menu use case
type Service struct {
	MenuRepository menuRepository.Repository
}

// GetAll is a function that returns all menus
func (s *Service) GetAll(page int64, limit int64) (*PaginationResultMenu, error) {
	all, err := s.MenuRepository.GetAll(page, limit)
	if err != nil {
		return nil, err
	}

	return &PaginationResultMenu{
		Data:       all.Data,
		Total:      all.Total,
		Limit:      all.Limit,
		Current:    all.Current,
		NextCursor: all.NextCursor,
		PrevCursor: all.PrevCursor,
		NumPages:   all.NumPages,
	}, nil
}

// GetByID is a function that returns a menu by id
func (s *Service) GetByID(id int) (*menuDomain.Menu, error) {
	return s.MenuRepository.GetByID(id)
}

// Create is a function that creates a menu
func (s *Service) Create(menu *NewMenu) (*menuDomain.Menu, error) {
	menuModel := menu.toDomainMapper()
	return s.MenuRepository.Create(menuModel)
}

// Delete is a function that deletes a menu by id
func (s *Service) Delete(id int) error {
	return s.MenuRepository.Delete(id)
}
