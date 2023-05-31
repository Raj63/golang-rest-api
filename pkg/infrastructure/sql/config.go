package sql

import "time"

// Config indicates database connection configuration.
type Config struct {
	// ConnMaxLifetime indicates the maximum amount of time a connection may be reused.
	ConnMaxLifetime time.Duration

	// ConnMaxIdleTime indicates the maximum amount of time a connection may be idle.
	ConnMaxIdleTime time.Duration

	// DriverName indicates the SQL driver to use, currently only supports:
	// 	- mysql
	// 	- postgres
	DriverName DriverName

	// MaxIdleConns indicates the maximum number of connections in the idle connection pool.
	MaxIdleConns int

	// MaxOpenConns indicates the maximum number of open connections to the database.
	MaxOpenConns int

	// SchemaSearchPath indicates the schema search path which is only used with "postgres".
	SchemaSearchPath string

	// SchemaMigrationsTable indicates the table name for storing the schema migration versions.
	SchemaMigrationsTable string

	// URI indicates the database connection string to connect.
	//
	// URI connection string documentation:
	//   - mysql: https://dev.mysql.com/doc/refman/8.0/en/connecting-using-uri-or-key-value-pairs.html#connecting-using-uri
	//   - postgres: https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	URI string
}

func defaultConfig(c *Config) *Config {
	if c.ConnMaxIdleTime == 0 {
		c.ConnMaxIdleTime = 5 * time.Minute
	}

	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 5 * time.Minute
	}

	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 16
	}

	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 16
	}

	if c.MaxIdleConns > c.MaxOpenConns {
		c.MaxIdleConns = c.MaxOpenConns
	}

	if c.SchemaSearchPath == "" {
		c.SchemaSearchPath = "public"
	}

	if c.SchemaMigrationsTable == "" {
		c.SchemaMigrationsTable = "schema_migrations"
	}

	return c
}
