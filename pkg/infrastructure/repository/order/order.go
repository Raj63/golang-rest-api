// Package order contains the repository implementation for the order entity
package order

import (
	"context"
	"database/sql"
	"errors"

	appErr "github.com/Raj63/golang-rest-api/pkg/domain/errors"
	domainOrder "github.com/Raj63/golang-rest-api/pkg/domain/order"
	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	sdksql "github.com/Raj63/golang-rest-api/pkg/infrastructure/sql"
	"github.com/go-sql-driver/mysql"
)

// Repository is a struct that contains the database implementation for order entity
type Repository struct {
	Store  *sdksql.DB
	Logger *logger.Logger
}

// Create ... Insert New data
func (r *Repository) Create(ctx context.Context, newOrder *domainOrder.Request) (*domainOrder.Request, error) {
	order := fromDomainMapper(newOrder)
	// store into DB
	tx, err := r.Store.DB().Beginx()
	if err != nil {
		return nil, err
	}
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	result, err := tx.NamedExecContext(ctx, "INSERT INTO orders (diner_id, menu_id, quantity, created_at, updated_at) VALUES (:diner_id, :menu_id, :quantity, NOW(), NOW());", order)
	if err != nil {
		_ = tx.Rollback()
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, appErr.NewAppErrorWithType(appErr.ResourceAlreadyExists)
		}
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
	return order.toDomainMapper(), nil
}

// GetByID ... Fetch only one order by Id
func (r *Repository) GetByID(ctx context.Context, dinerID int64) ([]domainOrder.Response, error) {
	var orders []Response

	err := r.Store.DB().Select(&orders, `
	SELECT 
		o.id,
		d.name as diner_name,
		m.name as menu_name, 
		m.description as menu_description,
		o.quantity,
		o.created_at 
	FROM orders o 
	INNER JOIN menus m 
		ON o.menu_id = m.id
	INNER JOIN diners d
		ON o.diner_id = d.id
	WHERE d.id = ?;`, dinerID)
	if err != nil {
		return nil, err
	}

	return arrayToDomainMapper(&orders), nil
}

// Delete ... Delete order
func (r *Repository) Delete(ctx context.Context, id int) (err error) {
	result, err := r.Store.DB().ExecContext(ctx, `
	DELETE FROM
		orders
	WHERE id = ?;
	`, id)
	if err != nil {
		r.Logger.ErrorfContext(ctx, "error when executing query. error: %+v, param: %+v\n", err, id)
		return err
	}

	rowAffected, errGetAffectedRow := result.RowsAffected()
	if errGetAffectedRow != nil || rowAffected == 0 {
		r.Logger.Errorf("error when get affected row. error: %+v", errGetAffectedRow)
		return appErr.NewAppErrorWithType(appErr.NotFound)
	}

	return nil
}
