package command

import (
	"fmt"

	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/mpon/ecswalk/internal/pkg/fuzzyfinder"
	"github.com/spf13/cobra"
)

// NewCmdServices represents services cmd
func NewCmdServices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "describe ECS services by selecting cluster interactively",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := awsapi.NewClient()
			if err != nil {
				return err
			}
			clusters, err := client.GetAllEcsClusters()
			if err != nil {
				return err
			}

			if len(clusters) == 0 {
				fmt.Println("cluster not found")
				return nil
			}

			cluster, err := fuzzyfinder.FindCluster(clusters)
			if err != nil {
				// Abort fuzzyfinder
				return nil
			}

			if err := runGetServices(client, cluster); err != nil {
				return nil
			}

			return nil
		},
	}
	return cmd
}
