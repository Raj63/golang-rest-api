// Package menu provides the use case for menu
package menu

import (
	domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"
)

func (n *NewMenu) toDomainMapper() *domainMenu.Menu {
	return &domainMenu.Menu{
		Name:        n.Name,
		Description: n.Description,
	}
}
