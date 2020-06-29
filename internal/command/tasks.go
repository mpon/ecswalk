package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
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

			describeServicesOutputs, err := client.DescribeAllECSServices(cluster)
			if err != nil {
				return err
			}

			services := []ecs.Service{}
			for _, describeServiceOutput := range describeServicesOutputs {
				for _, service := range describeServiceOutput.Services {
					services = append(services, service)
				}
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

			getTasksCmdRun(*cluster.ClusterName, *service.ServiceName)
			return nil
		},
	}
	return cmd
}
