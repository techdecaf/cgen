package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/techdecaf/cgen/app"
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Creates a new git tag with an increase in the current semversion i.e. v1.0.2",
	Long:  `Creates a new git tag with an increase in the current semversion i.e. v1.0.2`,
	Run: func(cmd *cobra.Command, args []string) {
		// parse flags
		var level, pattern string
		var err error

		if level, err = cmd.Flags().GetString("level"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}
		if pattern, err = cmd.Flags().GetString("pattern"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}

		// initialize a new instance of cgen
		cgen := &app.CGen{}
		if err := cgen.Init(); err != nil {
			app.Log.Fatal("cgen_init", err)
		}

		ver, err := app.Bump(level, pattern)
		if err != nil {
			app.Log.Fatal("app_bump", err)
		}

		fmt.Println(ver)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(bumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	bumpCmd.Flags().StringP("level", "l", "patch", "accepts (major, minor, patch or pre-release) strings")
	bumpCmd.Flags().StringP("pattern", "c", "v%n", "use a custom pattern for the git tag, defaults to v%s i.e. (v1.0.2)")
}
