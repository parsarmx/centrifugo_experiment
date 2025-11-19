package cmd

// import (
// 	"golang_template/internal/database/postgres"
// 	"log"

// 	"github.com/spf13/cobra"
// )

// // postgresMakemigrationCmd represents the pg_makemigration command
// var postgresMakemigrationCmd = &cobra.Command{
// 	Use:   "pg_makemigration [name]",
// 	Short: "Create a new migration file",
// 	Long: `Create a new migration file. Example:
// 		pg_makemigration add_users_table // Create a new migration file with sql type
// 		pg_makemigration add_users_table --type go // Create a new migration file with go type
// 		pg_makemigration add_users_table --dir ./database/migrations // Create a new migration file in specific directory`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		cmd.Println("pg_makemigration")
// 		name := args[0]
// 		if name == "" {
// 			cmd.PrintErrf("Migration name is required")
// 			return
// 		}

// 		typeFlag, err := cmd.Flags().GetString("type")
// 		if err != nil {
// 			cmd.PrintErrf("Error while getting type flag:\n\t %v", err)
// 			return
// 		}
// 		if typeFlag != "sql" && typeFlag != "go" {
// 			cmd.PrintErrf("Invalid migration type: %s. Must be 'sql' or 'go'", typeFlag)
// 			return
// 		}

// 		dirFlag, err := cmd.Flags().GetString("dir")
// 		if err != nil {
// 			cmd.PrintErrf("Error while getting dir flag:\n\t %v", err)
// 			return
// 		}

// 		// dbConfig, err := config.LoadConfig("config/config.yaml")

// 		if err != nil {
// 			log.Fatalf("failed to setup viper: %s", err.Error())
// 			return
// 		}

// 		// db, err := postgres.NewDatabase(cmd.Context(), &dbConfig.DB)
// 		if err != nil {
// 			cmd.PrintErrf("Error while initializing database:\n\t %v", err)
// 			return
// 		}
// 		// defer db.Close()

// 		migration := postgres.NewMigration(nil)
// 		err = migration.CreateFile(dirFlag, name, typeFlag)
// 		if err != nil {
// 			cmd.PrintErrf("Error while creating migration file:\n\t %v", err)
// 			return
// 		}
// 	},
// }

// func init() {
// 	RootCmd.AddCommand(postgresMakemigrationCmd)

// 	postgresMakemigrationCmd.Flags().String("dir", "./internal/database/postgres/migrations", "Directory of the migrations")
// 	postgresMakemigrationCmd.Flags().String("type", "sql", "Type of the migration file which is sql or go")
// }
