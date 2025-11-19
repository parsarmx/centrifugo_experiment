package cmd

import (
	"golang_template/internal/config"
	"golang_template/internal/database/postgres"
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// postgresMigrateCmd represents the pg_migrate command
var postgresMigrateCmd = &cobra.Command{
	Use:   "pg_migrate",
	Short: "Migrate the migration files",
	Long: `Migrate the migration files. Example:
	pg_migrate                              Migrate all the migration files
	pg_migrate --dir ./database/migrations  Migrate all the migration files from sepecific directory
	pg_migrate --version 1                  Migrate the migration file up to version 1
	pg_migrate --fake true                  Fake apply all the migration files`,
	Run: func(cmd *cobra.Command, args []string) {
		versionFlag, err := cmd.Flags().GetInt64("version")
		if err != nil {
			cmd.PrintErrf("Error while getting version flag:\n\t %v", err)
			return
		}

		dirFlag, err := cmd.Flags().GetString("dir")
		if err != nil {
			cmd.PrintErrf("Error while getting dir flag:\n\t %v", err)
			return
		}

		fakeFlag, err := cmd.Flags().GetBool("fake")
		if err != nil {
			cmd.PrintErrf("Error while getting fake flag:\n\t %v", err)
			return
		}

		dbConfig, err := config.LoadConfig("config/config.yaml")

		if err != nil {
			log.Fatalf("failed to setup viper: %s", err.Error())
			return
		}

		logger, err := zap.NewDevelopment()
		db, err := postgres.NewDatabase(cmd.Context(), &dbConfig.DB, logger)
		if err != nil {
			cmd.PrintErrf("Error while initializing database:\n\t %v", err)
			return
		}
		defer db.Close()

		migration := postgres.NewMigration(db.DB())
		err = migration.Apply(dirFlag, versionFlag, fakeFlag)
		if err != nil {
			cmd.PrintErrf("Error while applying migrations:\n\t %v", err)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(postgresMigrateCmd)
	postgresMigrateCmd.Flags().String("dir", "./internal/database/postgres/migrations", "Directory of the migrations")
	postgresMigrateCmd.Flags().Int64("version", 0, "Version of the migration that is going to be applied")
	postgresMigrateCmd.Flags().Bool("fake", false, "Fake apply migrations.")
}
