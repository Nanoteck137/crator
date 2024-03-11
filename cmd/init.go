package cmd

import (
	"log"
	"os"

	"github.com/nanoteck137/crator/app"
	"github.com/nanoteck137/crator/template"
	"github.com/nanoteck137/crator/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use: "init <TEMPLATE_NAME>",
	Short: "Initialize with specified template",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]
		output, _ := cmd.Flags().GetString("output")

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

		empty, err := utils.IsDirEmpty(output)

		if err == nil && !empty {
			log.Fatalf("'%s' is not empty", output)
		}

		err = templ.Execute(output)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	initCmd.Flags().StringP("output", "o", ".", "Output of create")

	rootCmd.AddCommand(initCmd)
}
