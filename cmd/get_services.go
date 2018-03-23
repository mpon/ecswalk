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

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/awsecs"
	"github.com/spf13/cobra"
)

var getServicesCmdFlagCluster string

// servicesCmd represents the services command
var getServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "get all ECS services specified cluster",
	Run: func(cmd *cobra.Command, args []string) {
		getServicesCmdRun(getServicesCmdFlagCluster)
	},
}

func getServicesCmdRun(cluster string) {
	describeServicesOutputs := awsecs.DescribeAllServices(cluster)
	services := []ecs.Service{}
	serviceArns := []string{}
	for _, describeServiceOutput := range describeServicesOutputs {
		for _, service := range describeServiceOutput.Services {
			services = append(services, service)
			serviceArns = append(serviceArns, *service.ServiceArn)
		}
	}
	describeTaskDefinitionOutputs := awsecs.DescribeTaskDefinitions(cluster, serviceArns)

	rows := GetServiceRows{}
	for _, describeTaskDefinitionOutput := range describeTaskDefinitionOutputs {
		taskDefinition := *describeTaskDefinitionOutput.TaskDefinition.TaskDefinitionArn
		service := awsecs.FindService(services, taskDefinition)

		for _, containerDefinition := range describeTaskDefinitionOutput.TaskDefinition.ContainerDefinitions {
			image, tag := awsecs.ShortDockerImage(*containerDefinition.Image)
			rows = append(rows, GetServiceRow{
				Name:           *service.ServiceName,
				TaskDefinition: awsecs.ShortArn(taskDefinition),
				Image:          image,
				Tag:            tag,
				DesiredCount:   *service.DesiredCount,
				RunningCount:   *service.RunningCount,
			})

		}
	}
	sort.Sort(rows)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprintln(w, "Name\tTaskDefinition\tImage\tTag\tDesired\tRunning\t")
	for _, row := range rows {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%d\t\n",
			row.Name,
			row.TaskDefinition,
			row.Image,
			row.Tag,
			row.DesiredCount,
			row.RunningCount)
	}
	w.Flush()
}

func init() {
	getCmd.AddCommand(getServicesCmd)
	getServicesCmd.Flags().StringVarP(&getServicesCmdFlagCluster, "cluster", "c", "", "AWS ECS cluster")
	getServicesCmd.MarkFlagRequired("cluster")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
