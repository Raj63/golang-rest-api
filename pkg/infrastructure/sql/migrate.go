package sql

import (
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	// this is needed to support mysql drivers
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	// this is needed to support postgres drivers
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

var (
	dbMigratePath = "db/migrate/primary"
)

type migrateDriver struct {
	config  *Config
	embedFS embed.FS
	httpfs.PartialDriver
}

// Open retrieves the FS handle based on the DB driver.
func (d *migrateDriver) Open(rawURL string) (source.Driver, error) {
	err := d.PartialDriver.Init(http.FS(d.embedFS), dbMigratePath)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Migrate is the DB migrator.
type Migrate struct {
	*migrate.Migrate
	migrateDriver
}

// NewMigrate initialises the DB migrator with open connection.
func NewMigrate(config *Config, embedFS embed.FS) (*Migrate, error) {
	md := migrateDriver{}
	md.config = config
	md.embedFS = embedFS
	source.Register("embed", &md)

	m, err := migrate.New("embed://", fmt.Sprintf("%s://%s", config.DriverName, config.URI))
	if err != nil {
		return nil, err
	}

	return &Migrate{m, md}, nil
}

// Run returns the migration status.
func (m *Migrate) Run() ([][]string, error) {
	status := [][]string{}
	entries, err := m.migrateDriver.embedFS.ReadDir(dbMigratePath)
	if err != nil {
		return nil, err
	}

	lastMigratedVersion, _, _ := m.Version()
	for _, entry := range entries {
		if !entry.IsDir() {
			fn := entry.Name()
			splits := strings.Split(fn, "_")
			if strings.HasSuffix(fn, ".up.sql") && len(splits) > 1 {
				version, err := strconv.Atoi(splits[0])
				if err != nil {
					return nil, err
				}

				state := "no"
				if version <= int(lastMigratedVersion) {
					state = "yes"
				}

				status = append(status, []string{fmt.Sprintf("%s/%s", dbMigratePath, fn), state})
			}
		}
	}

	return status, nil
}
