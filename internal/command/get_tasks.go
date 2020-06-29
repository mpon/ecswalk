package command

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/awsapi"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
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
			return getTasksCmdRun(getTasksCmdFlagCluster, getTasksCmdFlagService)
		},
	}
	cmd.Flags().StringVarP(&getTasksCmdFlagCluster, "cluster", "c", "", "AWS ECS cluster")
	cmd.MarkFlagRequired("cluster")
	cmd.Flags().StringVarP(&getTasksCmdFlagService, "service", "s", "", "AWS ECS service")
	cmd.MarkFlagRequired("service")

	return cmd
}

func getTasksCmdRun(clusterName string, service string) error {
	var cluster *ecs.Cluster
	client, err := awsapi.NewClient()
	if err != nil {
		return err
	}

	output, err := client.DescribeECSClusters()
	for _, c := range output.Clusters {
		c := c
		if *c.ClusterName == clusterName {
			cluster = &c
		}
	}

	o, err := client.DescribeECSServices(cluster, []string{service})
	if err != nil {
		return err
	}

	containerInstanceArns, rows, err := describeTasks(cluster, &o.Services[0])
	if err != nil {
		return err
	}
	instances := NewInstances(containerInstanceArns)

	ec2InstanceIds := []string{}

	if len(containerInstanceArns) > 0 {
		describeContainerInstancesOutput, err := client.DescribeContainerInstances(cluster, containerInstanceArns)
		if err != nil {
			return xerrors.Errorf("message: %w", err)
		}
		for _, containerInstance := range describeContainerInstancesOutput.ContainerInstances {
			instances.UpdateEC2InstanceIDByArn(*containerInstance.Ec2InstanceId, *containerInstance.ContainerInstanceArn)
			ec2InstanceIds = append(ec2InstanceIds, *containerInstance.Ec2InstanceId)
		}
		ec2InstanceIds = sliceutil.DistinctSlice(ec2InstanceIds)

		describeInstancesOutput, err := client.DescribeEC2Instances(ec2InstanceIds)
		if err != nil {
			return err
		}

		for _, reservation := range describeInstancesOutput.Reservations {
			for _, instance := range reservation.Instances {
				instances.UpdatePrivateIPByInstanceID(*instance.PrivateIpAddress, *instance.InstanceId)
			}
		}
	}

	for _, row := range rows {
		for _, data := range instances {
			if row.ContainerInstanceArn == data.ContainerInstanceArn {
				row.EC2InstanceID = data.EC2InstanceID
				row.PrivateIP = data.PrivateIP
			}
		}
	}

	sort.Sort(rows)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "TaskId\tTaskDefinition\tStatus\tEC2Instance\tPrivateIp")
	for _, row := range rows {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			row.TaskID,
			row.TaskDefinition,
			row.Status,
			row.EC2InstanceID,
			row.PrivateIP,
		)
	}
	w.Flush()
	return nil
}

func describeTasks(cluster *ecs.Cluster, service *ecs.Service) ([]string, GetTaskRows, error) {
	client, err := awsapi.NewClient()
	if err != nil {
		return []string{}, GetTaskRows{}, err
	}
	containerInstanceArns := []string{}
	rows := GetTaskRows{}

	listTasksOutput, err := client.ListECSTasks(cluster, service)
	if err != nil {
		return []string{}, GetTaskRows{}, err
	}
	var ecsTasks []ecs.Task = []ecs.Task{}

	if len(listTasksOutput.TaskArns) > 0 {
		describeTasksOutput, err := client.DescribeTasks(cluster, listTasksOutput.TaskArns)
		if err != nil {
			return []string{}, GetTaskRows{}, err
		}
		ecsTasks = describeTasksOutput.Tasks
	}

	for _, task := range ecsTasks {
		rows = append(rows, &GetTaskRow{
			TaskID:               awsapi.ShortArn(*task.TaskArn),
			TaskDefinition:       awsapi.ShortArn(*task.TaskDefinitionArn),
			Status:               *task.LastStatus,
			ContainerInstanceArn: *task.ContainerInstanceArn,
		})
		containerInstanceArns = append(containerInstanceArns, *task.ContainerInstanceArn)
	}
	containerInstanceArns = sliceutil.DistinctSlice(containerInstanceArns)

	return containerInstanceArns, rows, nil
}
