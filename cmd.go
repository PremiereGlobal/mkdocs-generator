package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

// Args contains any arguments to the program
var Args *viper.Viper

// rootCmd is the root-level command line command
var rootCmd = &cobra.Command{
	Use:   "mkdocs-generator",
	Short: "mkdocs-generator crawls Bitbucket and creates a directory structure to be consumed by mkdocs",
	Long: `Allows mkdocs site to be generated from an entire Bitbucket server.
                Mkdocs information can be found at https://www.mkdocs.org/`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		log.Fatal("Must call subcommand")
	},
}

// generateCmd is the command executed by the "generate" subcommand
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the directory structure",
	Long:  `Generates the directory structure to be consumed by mkdocs`,
	Run: func(cmd *cobra.Command, args []string) {
		generate()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `All software has versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mkdocs-generator/%s", version)
		fmt.Println("")
	},
}

func initCmd() {

	cobra.OnInitialize(initConfig)

	Args = viper.New()

	// All command line arguments can be set via environment variables in the form
	// of MG_<command line arg> with dashes replace by underscores.  For example,
	// MG_LOG_LEVEL=debug will set the log level
	Args.SetEnvPrefix("MG")
	Args.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	Args.SetEnvKeyReplacer(replacer)

	rootCmd.PersistentFlags().StringP("log-level", "v", "info", "Log Level (\"fatal\",\"warn\",\"info\",\"debug\")")
	Args.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	generateCmd.Flags().StringP("bitbucket-url", "b", "", "Bitbucket URL (ex: https://bitbucket.mycompany.com)")
	Args.BindPFlag("bitbucket-url", generateCmd.Flags().Lookup("bitbucket-url"))
	generateCmd.Flags().StringP("bitbucket-user", "u", "", "Bitbucket username")
	Args.BindPFlag("bitbucket-user", generateCmd.Flags().Lookup("bitbucket-user"))
	generateCmd.Flags().StringP("bitbucket-password", "p", "", "Bitbucket password")
	Args.BindPFlag("bitbucket-password", generateCmd.Flags().Lookup("bitbucket-password"))
	generateCmd.Flags().StringP("build-dir", "d", "build", "The directory to build out the markdown structure")
	Args.BindPFlag("build-dir", generateCmd.Flags().Lookup("build-dir"))
	generateCmd.Flags().StringP("docs-dir", "c", "", "Path an existing mkdocs structure (including a mkdocs.yml file)")
	Args.BindPFlag("docs-dir", generateCmd.Flags().Lookup("docs-dir"))
	generateCmd.Flags().StringP("mkdocs-key", "k", "Projects", "The mkdocs.yml path to write the project structure (--mkdocs-file must also be set)")
	Args.BindPFlag("mkdocs-key", generateCmd.Flags().Lookup("mkdocs-key"))

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	v := Args.GetString("log-level")
	switch v {
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		break
	default:
		log.Warn("Unknown log level ", v)
		log.SetLevel(logrus.InfoLevel)
	}
}

// Executes the root command
func executeCmd() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
