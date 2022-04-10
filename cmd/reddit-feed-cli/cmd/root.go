package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	addr    string
	pageID  uint64
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:               "reddit-feed-cli",
	Short:             "The command-line tool, reddit-feed-cli, allows you to communicate with the reddit-feed-api service.",
	DisableAutoGenTag: true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func init() {
	cobra.EnableCommandSorting = false

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(producerCmd)
	rootCmd.AddCommand(consumerCmd)

	rootCmd.AddCommand(versionCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
