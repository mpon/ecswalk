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

package awsecs

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecsctl/sliceutil"
)

// ListClusters to list clusters
func ListClusters() *ecs.ListClustersOutput {
	svc := newSvc()
	input := &ecs.ListClustersInput{}

	req := svc.ListClustersRequest(input)
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
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}

// DescribeClusters to describe a cluster
func DescribeClusters(clusterArns []string) *ecs.DescribeClustersOutput {
	svc := newSvc()
	input := &ecs.DescribeClustersInput{
		Clusters: clusterArns,
	}

	req := svc.DescribeClustersRequest(input)
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
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}
	return result
}

// ListServices to list ECS Service recursively
func ListServices(cluster string) []*ecs.ListServicesOutput {

	svc := newSvc()
	outputs, err := listServices(cluster, svc, nil, nil)

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
		return nil
	}

	return outputs
}

// DescribeServices to describe ECS services specified cluster
func DescribeServices(cluster string, services []string) *ecs.DescribeServicesOutput {

	svc := newSvc()

	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: services,
	}

	req := svc.DescribeServicesRequest(input)
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
		return nil
	}
	return result
}

// DescribeTaskDefinition describe specified task definition
func DescribeTaskDefinition(taskDefinitionArn string) []string {
	svc := newSvc()
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	}

	req := svc.DescribeTaskDefinitionRequest(input)
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
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	images := []string{}
	for _, c := range result.TaskDefinition.ContainerDefinitions {
		t := strings.Split(taskDefinitionArn, "/")
		s := strings.Split(*c.Image, "/")
		r := fmt.Sprintf("%s => %s", t[len(t)-1], s[len(s)-1])
		images = append(images, r)
	}
	return images
}

// DescribeTaskDefinitions describe with task definition about all services
func DescribeTaskDefinitions(cluster string, services []string) []string {
	maxAPILimitChunkSize := 10
	taskDefinitions := []string{}
	outputs := []string{}

	wg := &sync.WaitGroup{}
	for _, chunkedServices := range sliceutil.ChunkedSlice(services, maxAPILimitChunkSize) {
		wg.Add(1)
		go func(c []string) {
			defer wg.Done()
			describeServicesOutput := DescribeServices(cluster, c)
			for _, service := range describeServicesOutput.Services {
				taskDefinitions = append(taskDefinitions, *service.TaskDefinition)
			}
		}(chunkedServices)
	}
	wg.Wait()

	for _, t := range taskDefinitions {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			outputs = append(outputs, DescribeTaskDefinition(t)...)
		}(t)
	}
	wg.Wait()

	sort.Strings(outputs)
	return outputs
}

func newSvc() *ecs.ECS {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("failed to load config, " + err.Error())
	}
	return ecs.New(cfg)
}

func listServices(cluster string, svc *ecs.ECS, nextToken *string, outputs []*ecs.ListServicesOutput) ([]*ecs.ListServicesOutput, error) {
	input := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	if nextToken != nil {
		input = &ecs.ListServicesInput{
			Cluster:   aws.String(cluster),
			NextToken: nextToken,
		}
	}

	req := svc.ListServicesRequest(input)
	result, err := req.Send()

	if err != nil {
		return nil, err
	}

	outputs = append(outputs, result)

	if result.NextToken != nil {
		return listServices(cluster, svc, result.NextToken, outputs)
	}
	return outputs, nil
}
