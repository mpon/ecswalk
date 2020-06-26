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
			output, err := client.DescribeECSClusters()
			if err != nil {
				return err
			}

			if len(output.Clusters) == 0 {
				fmt.Println("cluster not found")
				return nil
			}

			cluster, err := fuzzyfinder.FindCluster(output.Clusters)
			if err != nil {
				// Abort fuzzyfinder
				return nil
			}

			getServicesCmdRun(*cluster.ClusterName)
			return nil
		},
	}
	return cmd
}
