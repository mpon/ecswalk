package command

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/spf13/cobra"
)

// NewCmdGetInstances represents the get instances command
func NewCmdGetInstances() *cobra.Command {
	var getServicesCmdFlagCluster string
	return &cobra.Command{
		Use:   "instances",
		Short: "get ECS container instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetInstancesCmd(getServicesCmdFlagCluster)
		},
	}
}

func runGetInstancesCmd(clusterName string) error {
	client, err := awsapi.NewClient()
	if err != nil {
		return err
	}
	cluster, err := client.GetECSCluster(clusterName)
	if err != nil {
		return err
	}
	return runGetInstances(client, cluster)
}

func runGetInstances(client *awsapi.Client, cluster *ecs.Cluster) error {
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
}
