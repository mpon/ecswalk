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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecsctl/awsecs"
	"github.com/spf13/cobra"
)

var describeServicesCmdFlagCluster string

// servicesCmd represents the services command
var describeServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "describe all ECS services specified cluster",
	Run: func(cmd *cobra.Command, args []string) {
		services := awsecs.ListServices(describeServicesCmdFlagCluster)

		chunked := [][]string{}
		maxAPILimitChunkSize := 10
		for i := 0; i < len(services); i += maxAPILimitChunkSize {
			end := i + maxAPILimitChunkSize

			if end > len(services) {
				end = len(services)
			}

			chunked = append(chunked, services[i:end])
		}
		svc := awsecs.NewSvc()
		for _, c := range chunked {
			p := describeParams{
				svc:      svc,
				services: c,
			}
			describeServices(p)
		}
	},
}

type describeParams struct {
	svc      *ecs.ECS
	services []string
}

func describeServices(p describeParams) {

	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(describeServicesCmdFlagCluster),
		Services: p.services,
	}

	req := p.svc.DescribeServicesRequest(input)
	result, err := req.Send()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	for _, s := range result.Services {
		fmt.Println(*s.TaskDefinition)
	}
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
