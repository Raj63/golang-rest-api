// Package menu contains the repository implementation for the menu entity
package menu

import domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"

func (menu *Menu) toDomainMapper() *domainMenu.Menu {
	return &domainMenu.Menu{
		ID:          menu.ID,
		Name:        menu.Name,
		Description: menu.Description,
		CreatedAt:   menu.CreatedAt,
	}
}

func fromDomainMapper(menu *domainMenu.Menu) *Menu {
	return &Menu{
		ID:          menu.ID,
		Name:        menu.Name,
		Description: menu.Description,
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
