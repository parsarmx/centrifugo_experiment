package cmd

// import (
// 	"golang_template/internal/config"
// 	"golang_template/internal/database/postgres"
// 	"log"

// 	"github.com/spf13/cobra"
// 	"go.uber.org/zap"
// )

// // postgresRollbackCmd represents the pg_rollback command
// var postgresRollbackCmd = &cobra.Command{
// 	Use:   "pg_rollback",
// 	Short: "Rollback migration/migrations using this command",
// 	Long: `Rollback migration file/files using goose. The path should be path to migration files.
// 		You can write the migration hash to rollback to a specific version.
// 		For example:
// 		pg_rollback
// 		pg_rollback --dir ./database/migrations --version 1`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		cmd.Println("pg_rollback")

// 		versionFlag, err := cmd.Flags().GetInt64("version")
// 		if err != nil {
// 			cmd.PrintErrf("Error while getting version flag:\n\t %v", err)
// 			return
// 		}

// 		dirFlag, err := cmd.Flags().GetString("dir")
// 		if err != nil {
// 			cmd.PrintErrf("Error while getting dir flag:\n\t %v", err)
// 			return
// 		}

// 		dbConfig, err := config.LoadConfig("config/config.yaml")

// 		if err != nil {
// 			log.Fatalf("failed to setup viper: %s", err.Error())
// 			return
// 		}

// 		logger, err := zap.NewDevelopment()
// 		db, err := postgres.NewDatabase(cmd.Context(), &dbConfig.DB, logger)
// 		if err != nil {
// 			cmd.PrintErrf("Error while initializing database:\n\t %v", err)
// 			return
// 		}
// 		defer db.Close()

// 		migration := postgres.NewMigration(db.DB())
// 		err = migration.Rollback(dirFlag, versionFlag)
// 		if err != nil {
// 			cmd.PrintErrf("Error while rolling back migrations:\n\t %v", err)
// 			return
// 		}
// 	},
// }

// func init() {
// 	RootCmd.AddCommand(postgresRollbackCmd)

// 	postgresRollbackCmd.Flags().String("dir", "./internal/database/postgres/migrations", "Directory of the migrations")
// 	postgresRollbackCmd.Flags().Int64("version", 0, "Version of the migration that migrations will be rolled back to")
// }
