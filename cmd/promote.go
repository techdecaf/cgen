package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	app "github.com/techdecaf/cgen/internal"
	"github.com/techdecaf/utils"
)

func doPromote(cmd *cobra.Command, args []string) {
	// parse flags
	var file []string
	var src, commit string
	var push, asTemplate, asPointer bool
	var err error

	if file, err = cmd.Flags().GetStringArray("files"); err != nil {
		app.Log.Fatal("cmd_flags", err)
	}

	if commit, err = cmd.Flags().GetString("commit"); err != nil {
		app.Log.Fatal("cmd_flags", err)
	}

	if asTemplate, err = cmd.Flags().GetBool("as-template"); err != nil {
		app.Log.Fatal("cmd_flags", err)
	}
	if asPointer, err = cmd.Flags().GetBool("as-pointer"); err != nil {
		app.Log.Fatal("cmd_flags", err)
	}

	if push, err = cmd.Flags().GetBool("push"); err != nil {
		app.Log.Fatal("cmd_flags", err)
	}

	if src, err = cmd.Flags().GetString("path"); err != nil {
		app.Log.Fatal("cmd_flags", err)
	}

	// initialize a new instance of cgen
	if err := cgen.Init(); err != nil {
		app.Log.Fatal("cgen_init", err)
	}

	params := app.GeneratorParams{
		ProjectDirectory: utils.PathTo(src), // destination directory for generated files
		PerformUpgrade:   false,             // run in upgrade mode
		PromoteFile:      true,              // run file promotion mode
		StaticOnly:       true,              // only copy static files, no template interpolation
		Verbose:          true,              // use verbose logging
	}

	if err := cgen.Generator.Init(params); err != nil {
		app.Log.Fatal("generator_init", err)
	}

	for _, file := range file {
		var source = path.Join(cgen.Generator.Project.Directory, file)
		var template = path.Join(cgen.Generator.Template.Files, file)

		if asTemplate {
			template = fmt.Sprintf("%v.tmpl", template)
		}
		if asPointer {
			var repo, filePath string
			var f *os.File
			template = fmt.Sprintf("%v.ptr", template)
			repoQuestion := app.Question{
				Name:   "repository",
				Type:   "string",
				Prompt: "Repository? ([user@]server:project.git format)",
			}
			filePathQuestion := app.Question{
				Name:   "file_path",
				Type:   "string",
				Prompt: "Repository file path? (path of the [user@]server:project.git:/template/test/jest.config.js is template/test)",
			}
			if repo, err = cgen.Generator.Ask(repoQuestion); err != nil {
				app.Log.Fatal("ptr_repository", err)
			}
			if filePath, err = cgen.Generator.Ask(filePathQuestion); err != nil {
				app.Log.Fatal("ptr_file_path", err)
			}
			if f, err = os.Create(source); err != nil {
				app.Log.Fatal("ptr_file_create", err)
			}
			fmt.Fprintln(f, fmt.Sprintf("repository=%s", repo))
			fmt.Fprintln(f, fmt.Sprintf("path=%s", filePath))
		}

		if err := cgen.Generator.Copy(source, template); err != nil {
			app.Log.Fatal("cgen_promote", err)
		}
	}

	if commit != "" {
		fmt.Printf("the --commit flag has not yet been implemented, please consider a pull request (%v) \n", commit)
	}

	if push {
		fmt.Printf("the --push flag has not yet been implemented, please consider a pull request (%v) \n", push)
	}

}

// annotations := []string{}
// annotation := make(map[string][]string)
// annotation[cobra.BashCompFilenameExt] = annotations

// promoteCmd represents the promote command
var promoteCmd = &cobra.Command{
	Use:   "promote",
	Short: "promote a file from a project to your cgen template",
	Long: `this command takes a file that you have modified in your current project
  and uses it to overrite the coresponding file in your cgen template.
  `,
	Run: doPromote,
}

func init() {
	rootCmd.AddCommand(promoteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// promoteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// promoteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	promoteCmd.Flags().StringArrayP("files", "f", []string{}, "relative file path to the files you wish to promote.")
	promoteCmd.MarkFlagFilename("files", "*")
	promoteCmd.Flags().StringP("commit", "c", "", "commit the promoted file to your cgen template.")
	promoteCmd.Flags().StringP("path", "p", pwd, "the root directory containing a .cgen.yaml file")

	promoteCmd.Flags().Bool("as-template", false, "append .tmpl to the end of the file name")
	promoteCmd.Flags().Bool("as-pointer", false, "append .ptr to the end of the file name")
	promoteCmd.Flags().Bool("push", false, "push changes to your cgen template to its remote.")
}
