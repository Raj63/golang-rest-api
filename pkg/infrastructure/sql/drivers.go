package sql

// DriverName indicates the name of the driver
type DriverName string

const (
	// MYSQL refers to the mysql driver
	MYSQL DriverName = "mysql"
	// POSTGRES refers to the postgres driver
	POSTGRES DriverName = "postgres"
)

var embedInit = "init.sql"
