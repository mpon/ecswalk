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

	"github.com/mpon/ecswalk/awsecs"
	"github.com/spf13/cobra"
)

var walkTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "describe ECS tasks by selecting cluster and service interactively",
	Run: func(cmd *cobra.Command, args []string) {
		listClustersOutput := awsecs.ListClusters()
		clusterNames := []string{}
		for _, clusterArn := range listClustersOutput.ClusterArns {
			clusterNames = append(clusterNames, awsecs.ShortArn(clusterArn))
		}

		prompt := newPrompt(clusterNames, "Select Cluster")
		_, cluster, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		describeServicesOutputs := awsecs.DescribeAllServices(cluster)
		serviceNames := []string{}
		for _, describeServiceOutput := range describeServicesOutputs {
			for _, service := range describeServiceOutput.Services {
				serviceNames = append(serviceNames, awsecs.ShortArn(*service.ServiceArn))
			}
		}

		prompt2 := newPrompt(serviceNames, "Select Service")
		_, service, err := prompt2.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		getTasksCmdRun(cluster, service)
	},
}

func init() {
	walkCmd.AddCommand(walkTasksCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
