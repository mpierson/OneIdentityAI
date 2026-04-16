package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of sped",
	Long:  `Print the version number of sped.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Synchronization Project EDitor v0.2")
	},
}
