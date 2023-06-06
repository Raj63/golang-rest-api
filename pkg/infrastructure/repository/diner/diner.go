// Package diner contains the repository implementation for the diner entity
package diner

import (
	"context"
	"database/sql"

	domainDiner "github.com/Raj63/golang-rest-api/pkg/domain/diner"
	"github.com/Raj63/golang-rest-api/pkg/domain/errors"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/repository"
	sdksql "github.com/Raj63/golang-rest-api/pkg/infrastructure/sql"
)

// Repository is a struct that contains the database implementation for diner entity
type Repository struct {
	Store  *sdksql.DB
	Logger *logger.Logger
}

// GetTotalCount Fetch total count of diners
func (r *Repository) GetTotalCount(ctx context.Context) (int64, error) {
	var total int64
	err := r.Store.DB().Get(&total, `SELECT count(id) FROM diners`)
	if err != nil {
		r.Logger.ErrorfContext(ctx, "error fetching total count: %v", err)
		return 0, err
	}
	return total, nil
}

// GetAll Fetch all diner data
func (r *Repository) GetAll(ctx context.Context, page int64, limit int64) (*repository.PaginationResultDiner, error) {
	var diners []Diner
	total, err := r.GetTotalCount(ctx)
	if err != nil {
		return nil, err
	}
	offset := (page - 1) * limit

	// Read diner from DB based on limit and offset
	rows, err := r.Store.DB().Queryx(`
SELECT
	id, name, table_no, created_at, updated_at 
FROM diners
ORDER BY
	name ASC
LIMIT ? OFFSET ?;`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var diner Diner
		err := rows.StructScan(&diner)
		if err != nil {
			r.Logger.ErrorfContext(ctx, "error when iterating rows. error %+v", err)
		}
		diners = append(diners, diner)
	}

	numPages := (total + limit - 1) / limit
	var nextCursor, prevCursor uint
	if page < numPages {
		nextCursor = uint(page + 1)
	}
	if page > 1 {
		prevCursor = uint(page - 1)
	}

	return &repository.PaginationResultDiner{
		Data:       arrayToDomainMapper(&diners),
		Total:      total,
		Limit:      limit,
		Current:    page,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		NumPages:   numPages,
	}, nil
}

// Create ... Insert New data
func (r *Repository) Create(ctx context.Context, newDiner *domainDiner.Diner) (*domainDiner.Diner, error) {
	diner := fromDomainMapper(newDiner)
	// store into DB
	tx, err := r.Store.DB().Beginx()
	if err != nil {
		return nil, err
	}
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	result, err := tx.NamedExecContext(ctx, "INSERT INTO diners (name, table_no, created_at, updated_at) VALUES (:name, :table_no, NOW(), NOW());", diner)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if affectedRows == 0 {
		_ = tx.Rollback()
		return nil, sql.ErrNoRows
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		r.Logger.ErrorfContext(ctx, "Commit failed: %v", err)
		return nil, err
	}
	return diner.toDomainMapper(), nil
}

// GetByID ... Fetch only one diner by Id
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainDiner.Diner, error) {
	var diner Diner

	err := r.Store.DB().Get(&diner, `SELECT id, name, table_no, created_at, updated_at FROM diners WHERE id=?;`, id)
	if err != nil {
		return nil, err
	}

	return diner.toDomainMapper(), nil
}

// Delete ... Delete diner
func (r *Repository) Delete(ctx context.Context, id int64) (err error) {
	result, err := r.Store.DB().ExecContext(ctx, `
	DELETE FROM
		diners
	WHERE id = ?;
	`, id)
	if err != nil {
		r.Logger.ErrorfContext(ctx, "error when executing query. error: %+v, param: %+v\n", err, id)
		return err
	}

	rowAffected, errGetAffectedRow := result.RowsAffected()
	if errGetAffectedRow != nil || rowAffected == 0 {
		r.Logger.Errorf("error when get affected row. error: %+v", errGetAffectedRow)
		return errors.NewAppErrorWithType(errors.NotFound)
	}

	return nil
}
