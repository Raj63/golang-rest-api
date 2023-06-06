package cmd

import (
	"embed"
	"os"
	"runtime"

	"github.com/Raj63/golang-rest-api/pkg/infrastructure/logger"
	sdksql "github.com/Raj63/golang-rest-api/pkg/infrastructure/sql"
	"github.com/spf13/cobra"
)

func dbUpdateCommand(di CommandDI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-db",
		Short: "Run the DB migration and Seeding scripts",
		Run: func(cmd *cobra.Command, args []string) {

			di.Logger.Infof(
				"* Go Version: %s, Env: %s",
				runtime.Version(),
				os.Getenv("APP_ENV")+".env",
			)

			DBMigrateUp(di.Logger, di.DB, di.EmbedFS)
			DBSeedUp(di.Logger, di.DB, di.EmbedFS)
			di.Logger.Infof("* The DB is migrated/Seeded succesfully.")
		},
	}

	return cmd
}

// DBMigrateUp is a command to migrate Database
func DBMigrateUp(logger *logger.Logger, db *sdksql.DB, embedFS embed.FS) {
	m, err := sdksql.NewMigrate(db.Config(), embedFS)
	if err != nil {
		logger.Fatal(err)
	}

	msg := "Migrating '%s' database schema...%s\n"
	migrations, err := m.Run()
	if err != nil {
		logger.Fatal(err)
	}

	if len(migrations) < 1 {
		logger.Infof(msg, "/db/migrate/primary", "NO MIGRATION FOUND")
		return
	}

	if migrations[len(migrations)-1][1] == "yes" {
		logger.Infof(msg, "/db/migrate/primary", "ALL VERSION(S) MIGRATED")
		return
	}

	if err := m.Up(); err != nil {
		logger.Fatal(err)
	}
	logger.Infof(msg, "/db/migrate/primary", "DONE")

	srcErr, dbErr := m.Close()
	if dbErr != nil {
		logger.Fatalln(dbErr)
	}

	if srcErr != nil {
		logger.Fatalln(srcErr)
	}
}

// DBSeedUp is used to seed database
func DBSeedUp(logger *logger.Logger, db *sdksql.DB, embedFS embed.FS) {
	msg := "Seeding '%s' database...%s\n"
	logger.Infof(msg, "/db/migrate/primary", "")

	s, err := sdksql.NewSeed(db.Config(), embedFS, logger)
	if err != nil {
		logger.Fatalln(err)
	}
	if s == nil {
		logger.Infof(msg, "/db/migrate/primary", " NO SEED FILES FOUND")
		return
	}

	if err := s.Run(); err != nil {
		logger.Fatalln(err)
	}

	logger.Infof(msg, "/db/migrate/primary", " DONE")

	if err := s.Close(); err != nil {
		logger.Fatalln(err)
	}
}
