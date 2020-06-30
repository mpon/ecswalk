package command

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

// NewCmdGetTasks represents the get tasks command
func NewCmdGetTasks() *cobra.Command {
	var getTasksCmdFlagCluster string
	var getTasksCmdFlagService string
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "get Tasks specified service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetTasksCmd(getTasksCmdFlagCluster, getTasksCmdFlagService)
		},
	}
	cmd.Flags().StringVarP(&getTasksCmdFlagCluster, "cluster", "c", "", "AWS ECS cluster")
	cmd.MarkFlagRequired("cluster")
	cmd.Flags().StringVarP(&getTasksCmdFlagService, "service", "s", "", "AWS ECS service")
	cmd.MarkFlagRequired("service")

	return cmd
}

func runGetTasksCmd(clusterName string, serviceName string) error {
	client, err := awsapi.NewClient()
	if err != nil {
		return err
	}

	cluster, err := client.GetECSCluster(clusterName)
	if err != nil {
		return err
	}

	service, err := client.GetECSService(cluster, serviceName)
	if err != nil {
		return err
	}

	if err := runGetTasks(client, cluster, service); err != nil {
		return err
	}
	return nil
}

func runGetTasks(client *awsapi.Client, cluster *ecs.Cluster, service *ecs.Service) error {

	tasks, err := client.GetECSTasks(cluster, service)
	if err != nil {
		return err
	}

	taskInfoList := awsapi.NewECSTaskInfoList(tasks)
	containerInstanceArns := taskInfoList.ContainerInstanceArns()

	if len(containerInstanceArns) > 0 {
		containerInstances, err := client.GetECSContainerInstances(cluster, containerInstanceArns)
		if err != nil {
			return xerrors.Errorf("GetECSContainerInstances: %w", err)
		}

		taskInfoList.SetContainerInstances(containerInstances)

		instances, err := client.GetEC2Instances(taskInfoList.EC2InstanceIds())
		if err != nil {
			return err
		}

		taskInfoList.SetEC2Instances(instances)
	}

	sort.Sort(taskInfoList)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "TaskId\tTaskDefinition\tStatus\tEC2Instance\tPrivateIp")
	for _, info := range taskInfoList {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			info.ShortTaskArn(),
			info.ShortTaskDefinitionArn(),
			*info.Task.LastStatus,
			*info.Instance.InstanceId,
			*info.Instance.PrivateIpAddress,
		)
	}
	w.Flush()
	return nil
}
