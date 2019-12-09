package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/techdecaf/cgen/app"
)

// Plugin structure for zsh
type Plugin struct {
	name string
	path string
}

func (plug *Plugin) script() string {
	return path.Join(plug.path, fmt.Sprintf("_%s", plug.name))
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates zsh completion scripts",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		zsh := os.Getenv("ZSH")
		if zsh == "" {
			app.Log.Fatal("completion", fmt.Errorf("could not find the ZSH environmental variable, is ZSH installed?"))
		}

		plugin := &Plugin{
			name: rootCmd.Name(),
			path: path.Join(zsh, "/completions"),
		}

		if _, err := os.Stat(plugin.path); os.IsNotExist(err) {
			os.MkdirAll(plugin.path, 0700)
		}

		if err := rootCmd.GenZshCompletionFile(plugin.script()); err != nil {
			app.Log.Fatal("completion", err)
		}

		fmt.Printf("a zsh completion file has been generated in %s \n", plugin.path)
		fmt.Println()
		fmt.Println("to utilize the plugin, please add 'compinit' to the end of your .zshrc file")

	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
