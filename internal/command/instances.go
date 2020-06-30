package command

import (
	"fmt"
	"os"
	"text/tabwriter"

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
			clusters, err := client.GetAllECSClusters()
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

			containerInstances, err := client.GetAllECSContainerInstances(cluster)
			if err != nil {
				return err
			}

			if len(containerInstances) == 0 {
				fmt.Printf("%s has no container instances\n", *cluster.ClusterName)
				return nil
			}

			cList := awsapi.NewECSContainerInstanceInfoList(containerInstances)

			ec2Instances, err := client.GetEC2Instances(cList.Ec2InstanceIds())
			if err != nil {
				return err
			}

			cList.SetEC2Instances(ec2Instances)

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 1, '\t', 0)
			fmt.Fprintln(w, "ContainerInstance\tEC2Instance\tPrivateIP\tConnected\tStatus\tRunning\tCPU\tMemory\tAgent\tDocker")
			for _, info := range cList {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\t%d\t%d\t%d\t%s\t%s\n",
					info.ShortContainerInstanceArn(),
					*info.Ec2Instance.InstanceId,
					*info.Ec2Instance.PrivateIpAddress,
					*info.ContainerInstance.AgentConnected,
					*info.ContainerInstance.Status,
					*info.ContainerInstance.RunningTasksCount,
					*info.CPUAvailable(),
					*info.MemoryAvailable(),
					*info.ContainerInstance.VersionInfo.AgentVersion,
					info.DockerVersion(),
				)
			}
			w.Flush()

			return nil
		},
	}
	return cmd
}
