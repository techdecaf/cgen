package cmd

import (
	"os"

	"github.com/spf13/cobra"
	app "github.com/techdecaf/cgen/internal"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a new generator from a git repository",
	Long:  `Installs a new generator from a git repository`,
	Run: func(cmd *cobra.Command, args []string) {
		// initialize a new instance of cgen
		cgen := &app.CGen{}
		if err := cgen.Init(); err != nil {
			app.Log.Fatal("cgen_init", err)
		}

		if _, err := cgen.Install(args[0]); err != nil {
			app.Log.Fatal("cgen_install", err)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
