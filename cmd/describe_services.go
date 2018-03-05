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
	"sync"

	"github.com/mpon/ecsctl/awsecs"
	"github.com/mpon/ecsctl/sliceutil"
	"github.com/spf13/cobra"
)

var describeServicesCmdFlagCluster string

// servicesCmd represents the services command
var describeServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "describe all ECS services specified cluster",
	Run: func(cmd *cobra.Command, args []string) {
		maxAPILimitChunkSize := 10
		taskDefinitions := []string{}
		outputs := []string{}

		services := awsecs.ListServices(describeServicesCmdFlagCluster)

		wg := &sync.WaitGroup{}
		for _, chunkedServices := range sliceutil.ChunkedSlice(services, maxAPILimitChunkSize) {
			wg.Add(1)
			go func(c []string) {
				defer wg.Done()
				ts := awsecs.DescribeServices(describeServicesCmdFlagCluster, c)
				taskDefinitions = append(taskDefinitions, ts...)
			}(chunkedServices)
		}
		wg.Wait()

		for _, t := range taskDefinitions {
			wg.Add(1)
			go func(t string) {
				defer wg.Done()
				outputs = append(outputs, awsecs.DescribeTaskDefinition(t)...)
			}(t)
		}
		wg.Wait()

		sort.Strings(outputs)

		for _, o := range outputs {
			fmt.Println(o)
		}

	},
}

func init() {
	describeCmd.AddCommand(describeServicesCmd)
	describeServicesCmd.Flags().StringVarP(&describeServicesCmdFlagCluster, "cluster", "c", "", "AWS ECS cluster)")
	describeServicesCmd.MarkFlagRequired("cluster")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
