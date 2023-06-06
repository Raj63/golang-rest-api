// Package diner contains the repository implementation for the diner entity
package diner

import domainDiner "github.com/Raj63/golang-rest-api/pkg/domain/diner"

func (diner *Diner) toDomainMapper() *domainDiner.Diner {
	return &domainDiner.Diner{
		ID:          diner.ID,
		Name:        diner.Name,
		TableNumber: diner.TableNumber,
		CreatedAt:   diner.CreatedAt,
		UpdatedAt:   diner.UpdatedAt,
	}
}

func fromDomainMapper(diner *domainDiner.Diner) *Diner {
	return &Diner{
		ID:          diner.ID,
		Name:        diner.Name,
		TableNumber: diner.TableNumber,
		CreatedAt:   diner.CreatedAt,
	}
}

func arrayToDomainMapper(diners *[]Diner) *[]domainDiner.Diner {
	dinersDomain := make([]domainDiner.Diner, len(*diners))
	for i, diner := range *diners {
		dinersDomain[i] = *diner.toDomainMapper()
	}

	return &dinersDomain
}
