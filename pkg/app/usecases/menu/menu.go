// Package menu provides the use case for menu
package menu

import (
	"context"

	menuDomain "github.com/Raj63/golang-rest-api/pkg/domain/menu"
	menuRepository "github.com/Raj63/golang-rest-api/pkg/infrastructure/repository/menu"
)

// Service is a struct that contains the repository implementation for menu use case
type Service struct {
	MenuRepository menuRepository.Repository
}

// GetAll is a function that returns all menus
func (s *Service) GetAll(ctx context.Context, page int64, limit int64) (*PaginationResultMenu, error) {
	all, err := s.MenuRepository.GetAll(ctx, page, limit)
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
func (s *Service) GetByID(ctx context.Context, id int64) (*menuDomain.Menu, error) {
	return s.MenuRepository.GetByID(ctx, id)
}

// GetByTopCount is a function that returns a menu by top counts
func (s *Service) GetByTopCount(ctx context.Context, count int) ([]menuDomain.Menu, error) {
	return s.MenuRepository.GetByTopCount(ctx, count)
}

// Create is a function that creates a menu
func (s *Service) Create(ctx context.Context, menu *NewMenu) (*menuDomain.Menu, error) {
	menuModel := menu.toDomainMapper()
	return s.MenuRepository.Create(ctx, menuModel)
}

// Delete is a function that deletes a menu by id
func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.MenuRepository.Delete(ctx, id)
}
