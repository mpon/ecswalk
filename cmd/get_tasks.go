// Copyright © 2018 Masato Oshima
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/mpon/ecsctl/awsec2"
	"github.com/mpon/ecsctl/awsecs"
	"github.com/mpon/ecsctl/sliceutil"
	"github.com/spf13/cobra"
)

var getTasksCmdFlagCluster string
var getTasksCmdFlagService string

// tasksCmd represents the tasks command
var getTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "get Tasks specified service",
	Run: func(cmd *cobra.Command, args []string) {
		rows := GetTaskRows{}
		listTasksOutput := awsecs.ListTasks(getTasksCmdFlagCluster, getTasksCmdFlagService)
		describeTasksOutput := awsecs.DescribeTasks(getTasksCmdFlagCluster, listTasksOutput.TaskArns)

		containerInstanceArns := []string{}
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
		containerInstances := []*ContainerInstance{}
		for _, arn := range containerInstanceArns {
			containerInstances = append(containerInstances, &ContainerInstance{
				ContainerInstanceArn: arn,
			})
		}

		ec2InstanceIds := []string{}
		describeContainerInstancesOutput := awsecs.DescribeContainerInstances(getTasksCmdFlagCluster, containerInstanceArns)
		for _, containerInstance := range describeContainerInstancesOutput.ContainerInstances {
			ec2InstanceIds = append(ec2InstanceIds, *containerInstance.Ec2InstanceId)
			for _, c := range containerInstances {
				if c.ContainerInstanceArn == *containerInstance.ContainerInstanceArn {
					c.EC2InstanceID = *containerInstance.Ec2InstanceId
				}
			}
			ec2InstanceIds = append(ec2InstanceIds, *containerInstance.Ec2InstanceId)
		}
		ec2InstanceIds = sliceutil.DistinctSlice(ec2InstanceIds)

		describeInstancesOutput := awsec2.DescribeInstances(ec2InstanceIds)

		for _, reservation := range describeInstancesOutput.Reservations {
			for _, instance := range reservation.Instances {
				for _, c := range containerInstances {
					if c.EC2InstanceID == *instance.InstanceId {
						c.PrivateIP = *instance.PrivateIpAddress
					}
				}
			}
		}

		for _, row := range rows {
			for _, c := range containerInstances {
				if row.ContainerInstanceArn == c.ContainerInstanceArn {
					row.PrivateIP = c.PrivateIP
				}
			}
		}

		sort.Sort(rows)

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 1, '\t', 0)
		fmt.Fprintln(w, "TaskId\tTaskDefinition\tStatus\tPrivateIp\tAwsLogs")
		for _, row := range rows {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				row.TaskID,
				row.TaskDefinition,
				row.Status,
				row.PrivateIP,
				row.AwsLogs,
			)
		}
		w.Flush()
	},
}

type ContainerInstance struct {
	ContainerInstanceArn string
	EC2InstanceID        string
	PrivateIP            string
}

func init() {
	getCmd.AddCommand(getTasksCmd)
	getTasksCmd.Flags().StringVarP(&getTasksCmdFlagCluster, "cluster", "c", "", "AWS ECS cluster)")
	getTasksCmd.MarkFlagRequired("cluster")
	getTasksCmd.Flags().StringVarP(&getTasksCmdFlagService, "service", "s", "", "AWS ECS cluster)")
	getTasksCmd.MarkFlagRequired("service")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tasksCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tasksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
