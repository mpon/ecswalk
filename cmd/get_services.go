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
	"sort"
	"strings"

	"github.com/mpon/ecsctl/awsecs"
	"github.com/spf13/cobra"
)

var getServicesCmdFlagCluster string

// servicesCmd represents the services command
var getServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "get all ECS services specified cluster",
	Run: func(cmd *cobra.Command, args []string) {
		listServicesOutputs := awsecs.ListServices(getServicesCmdFlagCluster)
		serviceArns := []string{}
		for _, listServiceOutput := range listServicesOutputs {
			for _, serviceArn := range listServiceOutput.ServiceArns {
				serviceArns = append(serviceArns, serviceArn)
			}
		}
		describeTaskDefinitionOutputs := awsecs.DescribeTaskDefinitions(getServicesCmdFlagCluster, serviceArns)

		sortTaskDefinitions := []string{}
		for _, describeTaskDefinitionOutput := range describeTaskDefinitionOutputs {
			names := strings.Split(*describeTaskDefinitionOutput.TaskDefinition.TaskDefinitionArn, "/")
			for _, container := range describeTaskDefinitionOutput.TaskDefinition.ContainerDefinitions {
				tdName := names[len(names)-1]
				images := strings.Split(*container.Image, "/")
				image := images[len(images)-1]
				o := fmt.Sprintf("%s => %s", tdName, image)
				sortTaskDefinitions = append(sortTaskDefinitions, o)
			}
		}
		sort.Strings(sortTaskDefinitions)

		for _, t := range sortTaskDefinitions {
			fmt.Println(t)
		}
	},
}

func init() {
	getCmd.AddCommand(getServicesCmd)
	getServicesCmd.Flags().StringVarP(&getServicesCmdFlagCluster, "cluster", "c", "", "AWS ECS cluster)")
	getServicesCmd.MarkFlagRequired("cluster")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
