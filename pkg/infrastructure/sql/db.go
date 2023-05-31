package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"

	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	"github.com/XSAM/otelsql"
	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zapadapter"
	"go.uber.org/zap"
)

// DB is a wrapper around sqlx.DB
type DB struct {
	master *sqlx.DB
	config *Config
	logger *logger.Logger
}

var (
	logLevel = sqldblogger.WithMinimumLevel(sqldblogger.LevelDebug)
)

func init() {
	if os.Getenv("APP_ENV") != "" && os.Getenv("APP_ENV") != "development" {
		logLevel = sqldblogger.WithMinimumLevel(sqldblogger.LevelInfo)
	}
}

// NewDB initialises the DB handle that is used to connect to a database.
func NewDB(c *Config, logger *logger.Logger) *DB {
	return &DB{
		config: defaultConfig(c),
		logger: logger,
	}
}

// DB returns the database.
func (db *DB) DB() *sqlx.DB {
	return db.master
}

// Open opens a database connection and verify with a ping.
//
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to
// close a DB.
func (db *DB) Open() error {
	driverName, err := otelsql.Register(string(db.config.DriverName))
	if err != nil {
		return err
	}

	masterDB, err := sqlx.Open(driverName, db.config.URI)
	if err != nil {
		return fmt.Errorf("unable to connect to '%s', error: %s", db.config.URI, err.Error())
	}

	wrappedMasterDB, err := newDBWithLogger(masterDB.DriverName(), db.config.URI, db.config, masterDB.Driver(), db.logger.Desugar())
	if err != nil {
		return err
	}
	db.master = wrappedMasterDB

	return nil
}

// Close returns the connection to the connection pool.
func (db *DB) Close() error {
	err := db.master.Close()
	if err != nil {
		return err
	}

	return nil
}

// Config returns the db configurations.
func (db *DB) Config() *Config {
	return db.config
}

// Conn return the database connection
func (db *DB) Conn(ctx context.Context) (*sql.Conn, error) {
	// get master connections
	masterConn, err := db.master.Conn(ctx)
	if err != nil {
		return nil, err
	}

	return masterConn, nil
}

func newDBWithLogger(driverName, uri string, config *Config, driver driver.Driver, logger *zap.Logger) (*sqlx.DB, error) {
	dbWithLogger := sqlx.NewDb(
		sqldblogger.OpenDriver(
			uri,
			driver,
			zapadapter.New(logger),
			sqldblogger.WithIncludeStartTime(true),
			sqldblogger.WithTimeFormat(sqldblogger.TimeFormatRFC3339),
			logLevel,
		),
		driverName,
	)
	if dbWithLogger == nil {
		return nil, fmt.Errorf("unable to create sqlx.DB with logger for '%s'", uri)
	}

	dbWithLogger.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	dbWithLogger.SetConnMaxLifetime(config.ConnMaxLifetime)
	dbWithLogger.SetMaxIdleConns(config.MaxIdleConns)
	dbWithLogger.SetMaxOpenConns(config.MaxOpenConns)
	if err := dbWithLogger.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping '%s', error: %s", uri, err.Error())
	}

	return dbWithLogger, nil
}

// Ping verifies the connections to the master/replica databases are still alive, establishing a
// connection if necessary.
//
// Ping uses context.Background internally; to specify the context, use PingContext.
func (db *DB) Ping() error {
	err := db.master.Ping()
	if err != nil {
		return fmt.Errorf("unable to ping '%s', error: %s", db.config.URI, err.Error())
	}

	return nil
}
