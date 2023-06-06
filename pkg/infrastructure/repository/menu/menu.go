// Package menu contains the repository implementation for the menu entity
package menu

import (
	"context"
	"database/sql"

	"github.com/Raj63/golang-rest-api/pkg/domain/errors"
	domainMenu "github.com/Raj63/golang-rest-api/pkg/domain/menu"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/repository"
	sdksql "github.com/Raj63/golang-rest-api/pkg/infrastructure/sql"
)

// Repository is a struct that contains the database implementation for menu entity
type Repository struct {
	Store  *sdksql.DB
	Logger *logger.Logger
}

// GetTotalCount Fetch total menu count
func (r *Repository) GetTotalCount(ctx context.Context) (int64, error) {
	var total int64
	err := r.Store.DB().Get(&total, `SELECT count(id) FROM menus`)
	if err != nil {
		r.Logger.ErrorfContext(ctx, "error fetching total count: %v", err)
		return 0, err
	}
	return total, nil
}

// GetAll Fetch all menu data
func (r *Repository) GetAll(ctx context.Context, page int64, limit int64) (*repository.PaginationResultMenu, error) {
	var menus []Menu
	total, err := r.GetTotalCount(ctx)
	if err != nil {
		return nil, err
	}
	offset := (page - 1) * limit

	// Read menu from DB based on limit and offset
	rows, err := r.Store.DB().Queryx(`
SELECT
	id, name, description, price, created_at, updated_at 
FROM menus
ORDER BY
	name ASC
LIMIT ? OFFSET ?;`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var menu Menu
		err := rows.StructScan(&menu)
		if err != nil {
			r.Logger.ErrorfContext(ctx, "error when iterating rows. error %+v", err)
		}
		menus = append(menus, menu)
	}

	numPages := (total + limit - 1) / limit
	var nextCursor, prevCursor uint
	if page < numPages {
		nextCursor = uint(page + 1)
	}
	if page > 1 {
		prevCursor = uint(page - 1)
	}

	return &repository.PaginationResultMenu{
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
func (r *Repository) Create(ctx context.Context, newMenu *domainMenu.Menu) (*domainMenu.Menu, error) {
	menu := fromDomainMapper(newMenu)
	// store into DB
	tx, err := r.Store.DB().Beginx()
	if err != nil {
		return nil, err
	}
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	result, err := tx.NamedExecContext(ctx, "INSERT INTO menus (name, description, price, created_at, updated_at) VALUES (:name, :description, :price, NOW(), NOW());", menu)
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
	return menu.toDomainMapper(), nil
}

// GetByID ... Fetch only one menu by Id
func (r *Repository) GetByID(ctx context.Context, id int64) (*domainMenu.Menu, error) {
	var menu Menu

	err := r.Store.DB().Get(&menu, `SELECT id, name, description, price, created_at, updated_at FROM menus WHERE id = ?;`, id)
	if err != nil {
		return nil, err
	}

	return menu.toDomainMapper(), nil
}

// GetByTopCount ... Fetch only top menus by count
func (r *Repository) GetByTopCount(ctx context.Context, count int) ([]domainMenu.Menu, error) {
	var menus []Menu

	err := r.Store.DB().Select(&menus, `SELECT m.id, m.name, m.description, m.price, SUM(o.quantity) AS count
	FROM menus m
	JOIN orders o ON m.id = o.menu_id
	GROUP BY m.id, m.name
	ORDER BY count DESC
	LIMIT ?;`, count)
	if err != nil {
		return nil, err
	}

	return *arrayToDomainMapper(&menus), nil
}

// Delete ... Delete menu
func (r *Repository) Delete(ctx context.Context, id int64) (err error) {
	result, err := r.Store.DB().ExecContext(ctx, `
	DELETE FROM
		menus
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
