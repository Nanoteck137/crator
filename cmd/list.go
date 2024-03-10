package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/nanoteck137/crator/app"
	"github.com/nanoteck137/crator/template"
	"github.com/spf13/cobra"
)


var listCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := app.ReadConfig()
		if err != nil {
			log.Fatal(err)
		}

		err = os.MkdirAll(config.Templates, 0755)
		if err != nil {
			log.Fatal(err)
		}

		availableTemplates, err := template.GetAvailable(config)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Available Templates:")
		for _, templ := range availableTemplates {
			fmt.Printf(" - %s\n", templ.Config.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
