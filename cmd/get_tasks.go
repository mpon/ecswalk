package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/mpon/ecswalk/internal/pkg/awsec2"
	"github.com/mpon/ecswalk/internal/pkg/awsecs"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
	"github.com/spf13/cobra"
)

var getTasksCmdFlagCluster string
var getTasksCmdFlagService string

// tasksCmd represents the tasks command
var getTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "get Tasks specified service",
	Run: func(cmd *cobra.Command, args []string) {
		getTasksCmdRun(getTasksCmdFlagCluster, getTasksCmdFlagService)
	},
}

func getTasksCmdRun(cluster string, service string) {
	containerInstanceArns, rows := describeTasks(cluster, service)
	instanceDatas := NewInstanceDatas(containerInstanceArns)

	ec2InstanceIds := []string{}
	describeContainerInstancesOutput := awsecs.DescribeContainerInstances(cluster, containerInstanceArns)
	for _, containerInstance := range describeContainerInstancesOutput.ContainerInstances {
		instanceDatas.UpdateEC2InstanceIDByArn(*containerInstance.Ec2InstanceId, *containerInstance.ContainerInstanceArn)
		ec2InstanceIds = append(ec2InstanceIds, *containerInstance.Ec2InstanceId)
	}
	ec2InstanceIds = sliceutil.DistinctSlice(ec2InstanceIds)

	describeInstancesOutput := awsec2.DescribeInstances(ec2InstanceIds)

	for _, reservation := range describeInstancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			instanceDatas.UpdatePrivateIPByInstanceID(*instance.PrivateIpAddress, *instance.InstanceId)
		}
	}

	for _, row := range rows {
		for _, data := range instanceDatas {
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
}

func describeTasks(cluster string, service string) ([]string, GetTaskRows) {
	containerInstanceArns := []string{}
	rows := GetTaskRows{}

	listTasksOutput := awsecs.ListTasks(cluster, service)
	describeTasksOutput := awsecs.DescribeTasks(cluster, listTasksOutput.TaskArns)

	for _, task := range describeTasksOutput.Tasks {
		rows = append(rows, &GetTaskRow{
			TaskID:               awsecs.ShortArn(*task.TaskArn),
			TaskDefinition:       awsecs.ShortArn(*task.TaskDefinitionArn),
			Status:               *task.LastStatus,
			ContainerInstanceArn: *task.ContainerInstanceArn,
		})
		containerInstanceArns = append(containerInstanceArns, *task.ContainerInstanceArn)
	}
	containerInstanceArns = sliceutil.DistinctSlice(containerInstanceArns)

	return containerInstanceArns, rows
}

func init() {
	getCmd.AddCommand(getTasksCmd)
	getTasksCmd.Flags().StringVarP(&getTasksCmdFlagCluster, "cluster", "c", "", "AWS ECS cluster")
	getTasksCmd.MarkFlagRequired("cluster")
	getTasksCmd.Flags().StringVarP(&getTasksCmdFlagService, "service", "s", "", "AWS ECS service")
	getTasksCmd.MarkFlagRequired("service")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tasksCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tasksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
