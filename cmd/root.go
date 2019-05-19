package cmd

import (
  "fmt"
  "os"
  // "strings"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  // "github.com/davecgh/go-spew/spew"
)

var rootCmd = &cobra.Command{
  Use:   "mkdocs-generator",
  Short: "mkdocs-generator crawls Bitbucket and creates a directory structure to be consumed by mkdocs",
  Long: `Allows mkdocs site to be generated from an entire Bitbucket server.
                Mkdocs information can be found at https://www.mkdocs.org/`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

var v *viper.Viper

var generateCmd = &cobra.Command{
  Use:   "generate",
  Short: "Generates the directory structure",
  Long: `Generates the directory structure to be consumed by mkdocs`,
  Run: func(cmd *cobra.Command, args []string) {
    // fmt.Println("Print: " + strings.Join(args, " "))
  },
}

func Init() {
  v = viper.New()
  generateCmd.Flags().StringP("bitbucket-url", "b", "", "Bitbucket URL (ex: https://bitbucket.mycompany.com)")
  v.BindPFlag("bitbucket-url", generateCmd.Flags().Lookup("bitbucket-url"))
  generateCmd.Flags().StringP("bitbucket-user", "u", "", "Bitbucket username")
  v.BindPFlag("bitbucket-user", generateCmd.Flags().Lookup("bitbucket-user"))
  generateCmd.Flags().StringP("bitbucket-password", "p", "", "Bitbucket password")
  v.BindPFlag("bitbucket-password", generateCmd.Flags().Lookup("bitbucket-password"))
  rootCmd.AddCommand(generateCmd)
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func Lookup(arg string) string {
  // spew.Dump(arg)
  // if arg == "bitbucket-url" {
    return v.GetString(arg)
  // }
  // return ""
}
