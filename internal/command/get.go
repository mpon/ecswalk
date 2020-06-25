package command

import (
	"github.com/spf13/cobra"
)

// NewCmdGet represents the get command
func NewCmdGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}
