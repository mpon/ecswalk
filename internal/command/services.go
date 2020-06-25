package command

import (
	"fmt"

	"github.com/mpon/ecswalk/internal/pkg/awsapi"
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
			clusterNames := []string{}
			for _, cluster := range output.Clusters {
				clusterNames = append(clusterNames, *cluster.ClusterName)
			}

			prompt := newPrompt(clusterNames, "Select Cluster")
			_, cluster, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return err
			}
			getServicesCmdRun(cluster)
			return nil
		},
	}
	return cmd
}
