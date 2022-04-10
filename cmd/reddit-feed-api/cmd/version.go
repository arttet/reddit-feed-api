package cmd

import (
	"fmt"

	"github.com/arttet/reddit-feed-api/internal/config"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Reddit Feed API",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", config.Version)
		fmt.Println("Commit Hash:", config.CommitHash)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
