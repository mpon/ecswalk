package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmdVersion represents the version command
func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of ecswalk",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ecswalk version %s\n", Version)
		},
	}
	return cmd
}
