package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func initCmd() {

	cobra.OnInitialize(initConfig)

	Args = viper.New()

	rootCmd.PersistentFlags().StringP("log-level", "v", "info", "Log Level (\"fatal\",\"warn\",\"info\",\"debug\")")
	Args.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))

	generateCmd.Flags().StringP("bitbucket-url", "b", "", "Bitbucket URL (ex: https://bitbucket.mycompany.com)")
	Args.BindPFlag("bitbucket-url", generateCmd.Flags().Lookup("bitbucket-url"))
	generateCmd.Flags().StringP("bitbucket-user", "u", "", "Bitbucket username")
	Args.BindPFlag("bitbucket-user", generateCmd.Flags().Lookup("bitbucket-user"))
	generateCmd.Flags().StringP("bitbucket-password", "p", "", "Bitbucket password")
	Args.BindPFlag("bitbucket-password", generateCmd.Flags().Lookup("bitbucket-password"))
	generateCmd.Flags().StringP("build-dir", "d", "build/docs", "The directory to build out the markdown structure")
	Args.BindPFlag("build-dir", generateCmd.Flags().Lookup("build-dir"))

	rootCmd.AddCommand(generateCmd)
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
