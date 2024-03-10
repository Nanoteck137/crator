package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/nanoteck137/crator/app"
	"github.com/nanoteck137/crator/template"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use: "create <TEMPLATE_NAME>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]
		output, _ := cmd.Flags().GetString("output")

		fmt.Printf("templateName: %v\n", templateName)
		fmt.Printf("output: %v\n", output)


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

		var templ *template.Template

		for _, t := range availableTemplates {
			if t.Config.Name == templateName {
				templ = &t
			}
		}

		if templ == nil {
			log.Fatalf("No template with name '%s'", templateName)
		}

		err = templ.Execute(output)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	createCmd.Flags().StringP("output", "o", ".", "Output of create")

	rootCmd.AddCommand(createCmd)
}
