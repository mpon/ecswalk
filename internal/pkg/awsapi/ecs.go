package awsapi

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
	"golang.org/x/sync/errgroup"
)

// DescribeECSClusters to describe clusters
func (client Client) DescribeECSClusters() (*ecs.DescribeClustersOutput, error) {
	listInput := &ecs.ListClustersInput{}
	listReq := client.ECSClient.ListClustersRequest(listInput)
	listRes, err := listReq.Send(context.Background())
	if err != nil {
		return nil, err
	}

	input := &ecs.DescribeClustersInput{
		Clusters: listRes.ClusterArns,
	}

	req := client.ECSClient.DescribeClustersRequest(input)
	res, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return res.DescribeClustersOutput, nil
}

// DescribeAllECSServices to describe all ECS services specified cluster
func (client Client) DescribeAllECSServices(cluster string) ([]*ecs.DescribeServicesOutput, error) {
	outputs, err := client.listECSServicesRecursively(cluster)
	if err != nil {
		return nil, err
	}
	serviceArns := []string{}
	for _, o := range outputs {
		for _, arn := range o.ServiceArns {
			serviceArns = append(serviceArns, arn)
		}
	}

	const maxAPILimitChunkSize = 10
	describeServicesOutputs := []*ecs.DescribeServicesOutput{}

	eg, ctx := errgroup.WithContext(context.Background())
	for _, s := range sliceutil.ChunkedSlice(serviceArns, maxAPILimitChunkSize) {
		s := s
		eg.Go(func() error {
			describeServicesOutput, err := client.describeECSServices(cluster, s)

			select {
			case <-ctx.Done():
				return nil
			default:
				if err != nil {
					return err
				}
				describeServicesOutputs = append(describeServicesOutputs, describeServicesOutput)
				return nil
			}
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return describeServicesOutputs, nil
}

// DescribeTaskDefinition to describe specified task definition
func (client Client) DescribeTaskDefinition(taskDefinitionArn string) (*ecs.DescribeTaskDefinitionOutput, error) {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	}

	req := client.ECSClient.DescribeTaskDefinitionRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return result.DescribeTaskDefinitionOutput, nil
}

// DescribeTaskDefinitions describe with task definition about all services
func (client Client) DescribeTaskDefinitions(cluster string, services []string) ([]*ecs.DescribeTaskDefinitionOutput, error) {
	const maxAPILimitChunkSize = 10
	taskDefinitions := []string{}
	outputs := []*ecs.DescribeTaskDefinitionOutput{}

	eg, ctx := errgroup.WithContext(context.Background())
	for _, s := range sliceutil.ChunkedSlice(services, maxAPILimitChunkSize) {
		s := s
		eg.Go(func() error {
			describeServicesOutput, err := client.describeECSServices(cluster, s)

			select {
			case <-ctx.Done():
				return nil
			default:
				if err != nil {
					return err
				}
				for _, service := range describeServicesOutput.Services {
					taskDefinitions = append(taskDefinitions, *service.TaskDefinition)
				}
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	for _, t := range taskDefinitions {
		t := t
		eg.Go(func() error {
			output, err := client.DescribeTaskDefinition(t)

			select {
			case <-ctx.Done():
				return nil
			default:
				if err != nil {
					return err
				}
				outputs = append(outputs, output)
				return nil
			}
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return outputs, nil
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

func (client Client) listECSServicesRecursively(cluster string) ([]*ecs.ListServicesOutput, error) {
	outputs, err := client.listECSServices(cluster, nil, nil)

	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func (client Client) listECSServices(cluster string, nextToken *string, outputs []*ecs.ListServicesOutput) ([]*ecs.ListServicesOutput, error) {
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
		return client.listECSServices(cluster, result.NextToken, outputs)
	}
	return outputs, nil
}

func (client Client) describeECSServices(cluster string, services []string) (*ecs.DescribeServicesOutput, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: services,
	}

	req := client.ECSClient.DescribeServicesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return result.DescribeServicesOutput, nil
}
