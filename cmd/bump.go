package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	app "github.com/techdecaf/cgen/internal"
)

// bumpCmd represents the bump command
var bumpCmd = &cobra.Command{
	Use:   "bump",
	Short: "Creates a new git tag with an increase in the current semantic version i.e. v1.0.2",
	Long:  `Creates a new git tag with an increase in the current semantic version i.e. v1.0.2`,
	Run: func(cmd *cobra.Command, args []string) {
		// parse flags
		var level, pattern string
		var dryRun, push bool
		var err error

		if level, err = cmd.Flags().GetString("level"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}

		if pattern, err = cmd.Flags().GetString("pattern"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}

		if dryRun, err = cmd.Flags().GetBool("dry-run"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}

		if push, err = cmd.Flags().GetBool("push"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}

		// initialize a new instance of cgen
		cgen := &app.CGen{}
		if err := cgen.Init(); err != nil {
			app.Log.Fatal("cgen_init", err)
		}

		ver, err := app.Bump(app.BumpParams{
			Place:   level,
			Pattern: pattern,
			DryRun:  dryRun,
			GitPush: push,
		})

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
	bumpCmd.Flags().StringP("level", "l", "patch", "accepts (major, minor, patch or pre-release); defaults to `patch`")
	bumpCmd.Flags().StringP("pattern", "p", "v%s", "use a custom pattern for the git tag, %s will be replaced with the version number, v%s => v1.0.2")
	bumpCmd.Flags().Bool("push", false, "runs git push --follow-tags after version increment")
	bumpCmd.Flags().BoolP("dry-run", "d", false, "dry run only, do not run git tag")
}
