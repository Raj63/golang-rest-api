package sql

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	"github.com/jmoiron/sqlx"
)

var (
	dbSeedPath = "db/seed/primary"
)

// Seed is the DB seeder.
type Seed struct {
	config  *Config
	db      *sqlx.DB
	embedFS embed.FS
	files   []string
	logger  *logger.Logger
}

// NewSeed initialises the DB seeder with open connection.
func NewSeed(config *Config, embedFS embed.FS, logger *logger.Logger) (*Seed, error) {
	entries, err := embedFS.ReadDir(dbSeedPath)
	if err != nil {
		return nil, err
	}

	files := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") && entry.Name() != embedInit {
			files = append(files, fmt.Sprintf("%s/%s", dbSeedPath, entry.Name()))
		}
	}

	if len(files) < 1 {
		// no seed file found
		return nil, nil
	}

	sort.Strings(files)

	db, err := sqlx.Open(string(config.DriverName), config.URI)
	if err != nil {
		return nil, err
	}

	return &Seed{config, db, embedFS, files, logger}, nil
}

// Close closes the database and prevents new queries from starting.
func (s *Seed) Close() error {
	return s.db.Close()
}

// Run executes the seeding.
func (s *Seed) Run() error {
	for _, file := range s.files {
		b, err := s.embedFS.ReadFile(file)
		if err != nil {
			return err
		}

		if len(b) == 0 {
			continue
		}

		if _, err := s.db.Exec(string(b)); err != nil {
			s.logger.Errorf("error executing: %v", string(b))
			return err
		}
	}

	return nil
}
