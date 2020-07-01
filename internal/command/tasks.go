package command

import (
	"fmt"

	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/mpon/ecswalk/internal/pkg/fuzzyfinder"
	"github.com/spf13/cobra"
)

// NewCmdTasks represents tasks cmd
func NewCmdTasks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "describe ECS tasks by selecting cluster and service interactively",
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

			services, err := client.GetAllEcsServices(cluster)
			if err != nil {
				return err
			}

			if len(services) == 0 {
				fmt.Printf("%s has no services\n", *cluster.ClusterName)
				return nil
			}

			service, err := fuzzyfinder.FindService(services)
			if err != nil {
				// Abort fuzzyfinder
				return nil
			}

			runGetTasks(client, cluster, service)
			return nil
		},
	}
	return cmd
}
