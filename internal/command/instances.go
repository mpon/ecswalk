package command

import (
	"fmt"

	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/mpon/ecswalk/internal/pkg/fuzzyfinder"
	"github.com/spf13/cobra"
)

// NewCmdInstances represents instances cmd
func NewCmdInstances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instances",
		Short: "describe ECS container instances by selecting cluster interactively",
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

			return runGetInstances(client, cluster)
		},
	}
	return cmd
}
