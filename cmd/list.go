package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
	"github.com/techdecaf/cgen/app"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "prints a list of currently installed directories",
	Long:  `prints a list of currently installed directories`,
	Run: func(cmd *cobra.Command, args []string) {
		// initialize a new instance of cgen
		cgen := &app.CGen{}
		if err := cgen.Init(); err != nil {
			app.Log.Fatal("cgen_init", err)
		}

		// list all available generators
		generators, err := cgen.ListInstalled()
		if err != nil {
			app.Log.Fatal("list_generators", err)
		}

		for _, template := range generators {
			fmt.Println(path.Join(cgen.TemplatesDir, template))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
