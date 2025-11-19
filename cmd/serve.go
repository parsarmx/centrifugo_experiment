package cmd

import (
	"context"
	application "golang_template/app"
	"golang_template/internal/config"
	"log"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "Run Server",
	Run:   runserver,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func runserver(cmd *cobra.Command, args []string) {
	log.Println("serve")
	//viper
	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to setup viper: %s", err.Error())
	}
	application := application.NewApplication(context.TODO(), config)
	application.Setup()
}
