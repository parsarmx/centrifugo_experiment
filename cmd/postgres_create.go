package cmd

import (
	"database/sql"
	"fmt"
	"golang_template/internal/config"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
)

// postgresCreateCmd represents the postgresCreate command
var postgresCreateCmd = &cobra.Command{
	Use:   "pg_create",
	Short: "creates postgres database",
	Long: `pg_create creates postgres database based on env variables
	pg_create --db_name=mydb creates database with name mydb`,
	Run: func(cmd *cobra.Command, args []string) {
		dbConfig, err := config.LoadConfig("config/config.yaml")
		if err != nil {
			log.Panicf("failed to setup viper: %s", err)
			return
		}

		databaseUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s?sslmode=disable", dbConfig.DB.User, dbConfig.DB.Password, dbConfig.DB.Host, dbConfig.DB.Port)

		conn, err := sql.Open("pgx", databaseUrl)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
			return
		}
		defer conn.Close()

		_, createErr := conn.Exec(fmt.Sprintf("CREATE DATABASE %s", dbConfig.DB.DBName))

		if createErr != nil {
			if pgErr, ok := createErr.(*pgconn.PgError); ok && pgErr.Code == "42P04" {
				cmd.Printf("Database %s already exists: %s.\n", dbConfig.DB.DBName, pgErr.Message)
				return
			} else {
				log.Fatalf("Failed to create the database: %v\n", createErr)
				return
			}
		} else {
			cmd.Printf("Database %s created successfully.\n", dbConfig.DB.DBName)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(postgresCreateCmd)
	postgresCreateCmd.Flags().String("db_name", "", "Database name")
}
