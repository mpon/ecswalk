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
	"text/tabwriter"

	"github.com/mpon/ecsctl/awsecs"
	"github.com/spf13/cobra"
)

var getTasksCmdFlagCluster string
var getTasksCmdFlagService string

// tasksCmd represents the tasks command
var getTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "get Tasks specified service",
	Run: func(cmd *cobra.Command, args []string) {
		listTasksOutput := awsecs.ListTasks(getTasksCmdFlagCluster, getTasksCmdFlagService)
		describeTasksOutput := awsecs.DescribeTasks(getTasksCmdFlagCluster, listTasksOutput.TaskArns)

		// task.id, taskdef, Status, external-link, image, tag, loggroup, logstream
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 1, '\t', 0)
		fmt.Fprintln(w, "TaskDef\tStatus\tLogStream")
		for _, task := range describeTasksOutput.Tasks {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				awsecs.ShortArn(*task.TaskDefinitionArn),
				*task.LastStatus,
				awsecs.ShortArn(*task.TaskArn),
			)
		}
		w.Flush()
	},
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
