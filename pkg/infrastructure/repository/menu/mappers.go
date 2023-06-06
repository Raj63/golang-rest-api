// Package menu contains the repository implementation for the menu entity
package menu

import domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"

func (menu *Menu) toDomainMapper() *domainMenu.Menu {
	return &domainMenu.Menu{
		ID:          menu.ID,
		Name:        menu.Name,
		Description: menu.Description,
		Price:       float64(menu.Price) / 100, // We stored price as numeric value in MYSQL database hence to divide with 100 to get the actual decial points
		CreatedAt:   menu.CreatedAt,
		UpdatedAt:   menu.UpdatedAt,
	}
}

func fromDomainMapper(menu *domainMenu.Menu) *Menu {
	return &Menu{
		ID:          menu.ID,
		Name:        menu.Name,
		Description: menu.Description,
		Price:       int(menu.Price * 100), // We store price as numeric value in MYSQL database hence to multiply with 100 to keep the decimal points
		CreatedAt:   menu.CreatedAt,
	}
}

func arrayToDomainMapper(menus *[]Menu) *[]domainMenu.Menu {
	menusDomain := make([]domainMenu.Menu, len(*menus))
	for i, menu := range *menus {
		menusDomain[i] = *menu.toDomainMapper()
	}

	return &menusDomain
}
