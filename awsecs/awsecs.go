package awsecs

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/sliceutil"
	"github.com/spf13/viper"
)

// ListClusters to list clusters
func ListClusters() *ecs.ListClustersOutput {
	svc := newSvc()
	input := &ecs.ListClustersInput{}

	req := svc.ListClustersRequest(input)
	result, err := req.Send(context.Background())
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
		os.Exit(1)
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
	result, err := req.Send(context.Background())
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
		os.Exit(1)
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
		os.Exit(1)
	}

	return outputs
}

// DescribeServices to describe ECS services specified cluster and services
func DescribeServices(cluster string, services []string) *ecs.DescribeServicesOutput {

	svc := newSvc()

	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: services,
	}

	req := svc.DescribeServicesRequest(input)
	result, err := req.Send(context.Background())
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
		os.Exit(1)
	}
	return result
}

// DescribeAllServices to describe all ECS services specified cluster
func DescribeAllServices(cluster string) []*ecs.DescribeServicesOutput {
	listServiceOutputs := ListServices(cluster)
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
			describeServicesOutput := DescribeServices(cluster, c)
			describeServicesOutputs = append(describeServicesOutputs, describeServicesOutput)
		}(chunkedServices)
	}
	wg.Wait()
	return describeServicesOutputs
}

// DescribeTaskDefinition to describe specified task definition
func DescribeTaskDefinition(taskDefinitionArn string) *ecs.DescribeTaskDefinitionOutput {
	svc := newSvc()
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	}

	req := svc.DescribeTaskDefinitionRequest(input)
	result, err := req.Send(context.Background())
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
		os.Exit(1)
	}
	return result
}

// DescribeTaskDefinitions describe with task definition about all services
func DescribeTaskDefinitions(cluster string, services []string) []*ecs.DescribeTaskDefinitionOutput {
	const maxAPILimitChunkSize = 10
	taskDefinitions := []string{}
	outputs := []*ecs.DescribeTaskDefinitionOutput{}

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
			outputs = append(outputs, DescribeTaskDefinition(t))
		}(t)
	}
	wg.Wait()

	return outputs
}

// ListTasks to list specified cluster
func ListTasks(cluster string, service string) *ecs.ListTasksOutput {
	svc := newSvc()
	input := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	}

	req := svc.ListTasksRequest(input)
	result, err := req.Send(context.Background())
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
	return result
}

// DescribeTasks to describe specified cluster and tasks
func DescribeTasks(cluster string, tasks []string) *ecs.DescribeTasksOutput {
	svc := newSvc()
	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   tasks,
	}

	req := svc.DescribeTasksRequest(input)
	result, err := req.Send(context.Background())
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
		os.Exit(1)
	}
	return result
}

// DescribeContainerInstances to describe container instances
func DescribeContainerInstances(cluster string, containerInstances []string) *ecs.DescribeContainerInstancesOutput {
	svc := newSvc()
	input := &ecs.DescribeContainerInstancesInput{
		Cluster:            aws.String(cluster),
		ContainerInstances: containerInstances,
	}

	req := svc.DescribeContainerInstancesRequest(input)
	result, err := req.Send(context.Background())
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
		os.Exit(1)
	}
	return result
}

func newSvc() *ecs.ECS {
	if viper.IsSet("profile") {
		cfg, err := external.LoadDefaultAWSConfig(
			external.WithSharedConfigProfile(viper.GetString("profile")),
		)
		if err != nil {
			panic("failed to load config, " + err.Error())
		}
		return ecs.New(cfg)
	}
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
	result, err := req.Send(context.Background())

	if err != nil {
		return nil, err
	}

	outputs = append(outputs, result)

	if result.NextToken != nil {
		return listServices(cluster, svc, result.NextToken, outputs)
	}
	return outputs, nil
}
