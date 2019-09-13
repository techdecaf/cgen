package cmd

import (
	"github.com/spf13/cobra"
	"github.com/techdecaf/cgen/app"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "this features is not currently supported pull request?",
	Long:  `this features is not currently supported pull request?`,
	Run: func(cmd *cobra.Command, args []string) {
		var static bool
		var err error

		if static, err = cmd.Flags().GetBool("static-only"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}

		// initialize a new instance of cgen
		cgen := &app.CGen{}
		if err := cgen.Init(); err != nil {
			app.Log.Fatal("cgen_init", err)
		}

		params := app.GeneratorParams{
			Destination:    pwd,    // destination directory for generated files
			PerformUpgrade: false,  // perform upgrade
			StaticOnly:     static, // only copy static files, no template interpolation
			Verbose:        true,   // use verbose logging
		}

		if err := cgen.Generator.Init(params); err != nil {
			app.Log.Fatal("generator_init", err)
		}

		if err := cgen.Generator.Exec(); err != nil {
			app.Log.Fatal("generator_exec", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	upgradeCmd.Flags().BoolP("static-only", "s", false, "does not generate template files (most commonly used with update)")
}
