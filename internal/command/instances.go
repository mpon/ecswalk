package command

import (
	"fmt"
	"strings"

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

			output2, err := client.ListContainerInstances(*cluster.ClusterName)
			if err != nil {
				return err
			}

			output3, err := client.DescribeContainerInstances(*cluster.ClusterName, output2.ContainerInstanceArns)
			if err != nil {
				return err
			}

			for _, c := range output3.ContainerInstances {
				var cpuAvailable int64
				var memoryAvailable int64
				for _, r := range c.RemainingResources {
					if *r.Name == "CPU" {
						cpuAvailable = *r.IntegerValue
					}
					if *r.Name == "MEMORY" {
						memoryAvailable = *r.IntegerValue
					}
				}
				fmt.Println(
					awsapi.ShortArn(*c.ContainerInstanceArn),
					*c.Ec2InstanceId,
					*c.AgentConnected,
					*c.Status,
					*c.RunningTasksCount,
					cpuAvailable,
					memoryAvailable,
					*c.VersionInfo.AgentVersion,
					strings.Replace(*c.VersionInfo.DockerVersion, "DockerVersion: ", "", 1),
				)
			}

			return nil
		},
	}
	return cmd
}
