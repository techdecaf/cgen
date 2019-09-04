package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/techdecaf/cgen/app"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// VERSION is converted to the git tag at compile time using the make build command.
var VERSION string

// local variables
var cfgFile string
var pwd, _ = os.Getwd()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cgen",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var name, template string
		var version, static bool
		var err error

		// Print Version
		if version, _ = cmd.Flags().GetBool("version"); version {
			fmt.Println(VERSION)
			os.Exit(0)
		}

		if name, err = cmd.Flags().GetString("name"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}

		if template, err = cmd.Flags().GetString("template"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}

		if static, err = cmd.Flags().GetBool("static-only"); err != nil {
			app.Log.Fatal("cmd_flags", err)
		}

		app.Log.Info("debugging", fmt.Sprintf("project: %s, template: %s, static: %v", name, template, static))

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

		// PERFORM PROJECT GENERATION
		if template == "" {
			template, err = cgen.Generator.Ask(app.Question{
				Name:    "Template",
				Type:    "select",
				Prompt:  "Pick a template.",
				Options: generators,
			})
		}

		here, err := os.Stat(pwd)
		if err != nil {
			app.Log.Fatal("current_dir", err)
		}

		if name == "" {
			name, err = cgen.Generator.Ask(app.Question{
				Name:    "Name",
				Type:    "string",
				Prompt:  "What do you want to call your project ",
				Default: here.Name(),
			})
		}

		// // check to see if directory is dirty.
		if files, err := ioutil.ReadDir(pwd); err != nil {
			app.Log.Fatal("ioutil.ReadDir(pwd)", err)
		} else {
			if len(files) != 0 {
				app.Confirm("this directory is not empty, do you want to continue?")
			}
		}

		params := app.GeneratorParams{
			Name:           name,              // name of this project
			TemplatesDir:   cgen.TemplatesDir, // directory of all cgen templates
			Tempate:        template,          // selected cgen template
			Destination:    pwd,               // destination directory for generated files
			PerformUpgrade: false,             // perform upgrade
			StaticOnly:     false,             // only copy static files, no template interpolation
			Verbose:        true,              // use verbose logging
		}

		if err := cgen.Generator.Init(params); err != nil {
			app.Log.Fatal("generator_init", err)
		}

		if err := cgen.Generator.Exec(); err != nil {
			log.Fatal("generator_exec", err)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cgen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("version", "v", false, "prints the cgen version number")

	rootCmd.Flags().BoolP("static-only", "s", false, "does not generate template files (most commonly used with update)")

	rootCmd.Flags().StringP("name", "n", "", "what do you want to call your newly generated project?")
	rootCmd.Flags().StringP("template", "t", "", "specify a which template you would like to use.")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cgen" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cgen")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}