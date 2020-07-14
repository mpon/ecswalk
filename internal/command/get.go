package command

import (
	"github.com/spf13/cobra"
)

// NewCmdGet represents the get command
func NewCmdGet() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}
