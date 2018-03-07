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
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/mpon/ecsctl/awsecs"
	"github.com/spf13/cobra"
)

// walkCmd represents the walk command
var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "describe ECS services by selecting cluster interactively",
	Run: func(cmd *cobra.Command, args []string) {
		listClustersOutput := awsecs.ListClusters()
		clusterNames := []string{}
		for _, clusterArn := range listClustersOutput.ClusterArns {
			s := strings.Split(clusterArn, "/")
			clusterNames = append(clusterNames, s[len(s)-1])
		}

		prompt := newPrompt(clusterNames)
		_, cluster, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		services := awsecs.ListServices(cluster)
		taskDefinitions := awsecs.DescribeTaskDefinitions(cluster, services)

		for _, t := range taskDefinitions {
			fmt.Println(t)
		}
	},
}

func newPrompt(clusters []string) promptui.Select {
	searcher := func(input string, index int) bool {
		cluster := strings.ToLower(clusters[index])
		return strings.Contains(cluster, input)
	}

	return promptui.Select{
		Label:    "Select cluster",
		Items:    clusters,
		Size:     20,
		Searcher: searcher,
	}
}

func init() {
	rootCmd.AddCommand(walkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
