package awsapi

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
)

// ListECSClusters to list ECS clusters
func (client Client) ListECSClusters() (*ecs.ListClustersOutput, error) {
	input := &ecs.ListClustersInput{}
	req := client.ECSClient.ListClustersRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return result.ListClustersOutput, nil
}

// DescribeClusters to describe a cluster
func (client Client) DescribeClusters(clusterArns []string) *ecs.DescribeClustersOutput {
	input := &ecs.DescribeClustersInput{
		Clusters: clusterArns,
	}

	req := client.ECSClient.DescribeClustersRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeException:
				fmt.Println(ecs.ErrCodeException, aerr.Error())
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
		os.Exit(1)
	}
	return result.DescribeClustersOutput
}

// ListServices to list ECS Service recursively
func (client Client) ListServices(cluster string) []*ecs.ListServicesOutput {

	outputs, err := client.listServices(cluster, nil, nil)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeException:
				fmt.Println(ecs.ErrCodeException, aerr.Error())
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
		os.Exit(1)
	}

	return outputs
}

// DescribeServices to describe ECS services specified cluster and services
func (client Client) DescribeServices(cluster string, services []string) *ecs.DescribeServicesOutput {
	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: services,
	}

	req := client.ECSClient.DescribeServicesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeException:
				fmt.Println(ecs.ErrCodeException, aerr.Error())
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
		os.Exit(1)
	}
	return result.DescribeServicesOutput
}

// DescribeAllServices to describe all ECS services specified cluster
func (client Client) DescribeAllServices(cluster string) []*ecs.DescribeServicesOutput {
	listServiceOutputs := client.ListServices(cluster)
	serviceArns := []string{}
	for _, listServiceOutput := range listServiceOutputs {
		for _, serviceArn := range listServiceOutput.ServiceArns {
			serviceArns = append(serviceArns, serviceArn)
		}
	}

	const maxAPILimitChunkSize = 10
	describeServicesOutputs := []*ecs.DescribeServicesOutput{}

	wg := &sync.WaitGroup{}
	for _, chunkedServices := range sliceutil.ChunkedSlice(serviceArns, maxAPILimitChunkSize) {
		wg.Add(1)
		go func(c []string) {
			defer wg.Done()
			describeServicesOutput := client.DescribeServices(cluster, c)
			describeServicesOutputs = append(describeServicesOutputs, describeServicesOutput)
		}(chunkedServices)
	}
	wg.Wait()
	return describeServicesOutputs
}

// DescribeTaskDefinition to describe specified task definition
func (client Client) DescribeTaskDefinition(taskDefinitionArn string) *ecs.DescribeTaskDefinitionOutput {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	}

	req := client.ECSClient.DescribeTaskDefinitionRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeException:
				fmt.Println(ecs.ErrCodeException, aerr.Error())
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
		os.Exit(1)
	}
	return result.DescribeTaskDefinitionOutput
}

// DescribeTaskDefinitions describe with task definition about all services
func (client Client) DescribeTaskDefinitions(cluster string, services []string) []*ecs.DescribeTaskDefinitionOutput {
	const maxAPILimitChunkSize = 10
	taskDefinitions := []string{}
	outputs := []*ecs.DescribeTaskDefinitionOutput{}

	wg := &sync.WaitGroup{}
	for _, chunkedServices := range sliceutil.ChunkedSlice(services, maxAPILimitChunkSize) {
		wg.Add(1)
		go func(c []string) {
			defer wg.Done()
			describeServicesOutput := client.DescribeServices(cluster, c)
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
			outputs = append(outputs, client.DescribeTaskDefinition(t))
		}(t)
	}
	wg.Wait()

	return outputs
}

// ListTasks to list specified cluster
func (client Client) ListTasks(cluster string, service string) *ecs.ListTasksOutput {
	input := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	}

	req := client.ECSClient.ListTasksRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeException:
				fmt.Println(ecs.ErrCodeException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			case ecs.ErrCodeServiceNotFoundException:
				fmt.Println(ecs.ErrCodeServiceNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	return result.ListTasksOutput
}

// DescribeTasks to describe specified cluster and tasks
func (client Client) DescribeTasks(cluster string, tasks []string) *ecs.DescribeTasksOutput {
	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   tasks,
	}

	req := client.ECSClient.DescribeTasksRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeException:
				fmt.Println(ecs.ErrCodeException, aerr.Error())
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
		os.Exit(1)
	}
	return result.DescribeTasksOutput
}

// DescribeContainerInstances to describe container instances
func (client Client) DescribeContainerInstances(cluster string, containerInstances []string) *ecs.DescribeContainerInstancesOutput {
	input := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(cluster),
		ContainerInstances: containerInstances,
	}

	req := client.ECSClient.DescribeContainerInstancesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeException:
				fmt.Println(ecs.ErrCodeException, aerr.Error())
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
		os.Exit(1)
	}
	return result.DescribeContainerInstancesOutput
}

func (client Client) listServices(cluster string, nextToken *string, outputs []*ecs.ListServicesOutput) ([]*ecs.ListServicesOutput, error) {
	input := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	if nextToken != nil {
		input = &ecs.ListServicesInput{
			Cluster:   aws.String(cluster),
			NextToken: nextToken,
		}
	}

	req := client.ECSClient.ListServicesRequest(input)
	result, err := req.Send(context.Background())

	if err != nil {
		return nil, err
	}

	outputs = append(outputs, result.ListServicesOutput)

	if result.NextToken != nil {
		return client.listServices(cluster, result.NextToken, outputs)
	}
	return outputs, nil
}
