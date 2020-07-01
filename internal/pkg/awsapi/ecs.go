package awsapi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

// GetEcsCluster to get an ECS cluster
func (client Client) GetEcsCluster(clusterName string) (*ecs.Cluster, error) {
	output, err := client.describeEcsClusters()
	if err != nil {
		return nil, err
	}

	for _, c := range output.Clusters {
		c := c
		if *c.ClusterName == clusterName {
			return &c, nil
		}
	}
	return nil, xerrors.Errorf("Not found ECS Cluster %s", clusterName)
}

// GetEcsService to get an ECS Service
func (client Client) GetEcsService(cluster *ecs.Cluster, serviceName string) (*ecs.Service, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  cluster.ClusterName,
		Services: []string{serviceName},
	}

	req := client.EcsClient.DescribeServicesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return &result.DescribeServicesOutput.Services[0], nil
}

// GetAllEcsServices to get all ECS Services
func (client Client) GetAllEcsServices(cluster *ecs.Cluster) ([]ecs.Service, error) {
	outputs, err := client.describeAllEcsServices(cluster)
	if err != nil {
		return nil, err
	}

	// flatten list
	var services []ecs.Service
	for _, o := range outputs {
		for _, s := range o.Services {
			s := s
			services = append(services, s)
		}
	}
	return services, nil
}

// GetAllEcsClusters to get all ECS Clusters
func (client Client) GetAllEcsClusters() ([]ecs.Cluster, error) {
	output, err := client.describeEcsClusters()
	if err != nil {
		return nil, err
	}
	return output.Clusters, nil
}

// GetEcsTaskDefinitions to get ECS task definition list
func (client Client) GetEcsTaskDefinitions(cluster *ecs.Cluster, services []ecs.Service) ([]ecs.TaskDefinition, error) {
	outputs, err := client.describeTaskDefinitions(cluster, services)
	if err != nil {
		return nil, err
	}

	var taskDefinitions []ecs.TaskDefinition
	for _, o := range outputs {
		taskDefinitions = append(taskDefinitions, *o.TaskDefinition)
	}
	return taskDefinitions, nil
}

func (client Client) describeEcsClusters() (*ecs.DescribeClustersOutput, error) {
	listInput := &ecs.ListClustersInput{}
	listReq := client.EcsClient.ListClustersRequest(listInput)
	listRes, err := listReq.Send(context.Background())
	if err != nil {
		return nil, err
	}

	input := &ecs.DescribeClustersInput{
		Clusters: listRes.ClusterArns,
	}

	req := client.EcsClient.DescribeClustersRequest(input)
	res, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return res.DescribeClustersOutput, nil
}

func (client Client) describeAllEcsServices(cluster *ecs.Cluster) ([]*ecs.DescribeServicesOutput, error) {
	outputs, err := client.listEcsServicesRecursively(cluster)
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
			describeServicesOutput, err := client.describeEcsServices(cluster, s)

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

func (client Client) describeTaskDefinition(taskDefinitionArn string) (*ecs.DescribeTaskDefinitionOutput, error) {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionArn),
	}

	req := client.EcsClient.DescribeTaskDefinitionRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return result.DescribeTaskDefinitionOutput, nil
}

func (client Client) describeTaskDefinitions(cluster *ecs.Cluster, services []ecs.Service) ([]*ecs.DescribeTaskDefinitionOutput, error) {
	const maxAPILimitChunkSize = 10
	taskDefinitions := []string{}
	outputs := []*ecs.DescribeTaskDefinitionOutput{}

	for _, s := range services {
		taskDefinitions = append(taskDefinitions, *s.TaskDefinition)
	}

	eg, ctx := errgroup.WithContext(context.Background())
	for _, t := range taskDefinitions {
		t := t
		eg.Go(func() error {
			output, err := client.describeTaskDefinition(t)

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

// GetEcsTasks to get ECS tasks of specified cluster and service
func (client Client) GetEcsTasks(cluster *ecs.Cluster, service *ecs.Service) ([]ecs.Task, error) {
	input := &ecs.ListTasksInput{
		Cluster:     cluster.ClusterName,
		ServiceName: service.ServiceName,
	}

	req := client.EcsClient.ListTasksRequest(input)
	output, err := req.Send(context.Background())
	if err != nil {
		return nil, xerrors.Errorf("ECS ListTasks: %w", err)
	}

	if len(output.TaskArns) == 0 {
		return []ecs.Task{}, nil
	}

	res, err := client.describeTasks(cluster, output.TaskArns)
	if err != nil {
		return nil, err
	}

	return res.Tasks, nil
}

func (client Client) describeTasks(cluster *ecs.Cluster, tasks []string) (*ecs.DescribeTasksOutput, error) {
	input := &ecs.DescribeTasksInput{
		Cluster: cluster.ClusterName,
		Tasks:   tasks,
	}

	req := client.EcsClient.DescribeTasksRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, xerrors.Errorf("ECS DescribeTasks: %w", err)
	}
	return result.DescribeTasksOutput, nil
}

// GetAllEcsContainerInstances to get all ECS container instances
func (client Client) GetAllEcsContainerInstances(cluster *ecs.Cluster) ([]ecs.ContainerInstance, error) {
	input := &ecs.ListContainerInstancesInput{
		Cluster: cluster.ClusterName,
	}

	req := client.EcsClient.ListContainerInstancesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}

	if len(result.ListContainerInstancesOutput.ContainerInstanceArns) == 0 {
		return []ecs.ContainerInstance{}, nil
	}

	return client.GetEcsContainerInstances(cluster, result.ListContainerInstancesOutput.ContainerInstanceArns)
}

// GetEcsContainerInstances to get container instances
func (client Client) GetEcsContainerInstances(cluster *ecs.Cluster, containerInstanceArns []string) ([]ecs.ContainerInstance, error) {
	input := &ecs.DescribeContainerInstancesInput{
		Cluster:            cluster.ClusterName,
		ContainerInstances: containerInstanceArns,
	}

	req := client.EcsClient.DescribeContainerInstancesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, xerrors.Errorf("ECS DescribeContainerInstances: %w", err)
	}
	return result.DescribeContainerInstancesOutput.ContainerInstances, nil
}

func (client Client) listEcsServicesRecursively(cluster *ecs.Cluster) ([]*ecs.ListServicesOutput, error) {
	outputs, err := client.listEcsServices(cluster, nil, nil)

	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func (client Client) listEcsServices(cluster *ecs.Cluster, nextToken *string, outputs []*ecs.ListServicesOutput) ([]*ecs.ListServicesOutput, error) {
	input := &ecs.ListServicesInput{
		Cluster: cluster.ClusterName,
	}

	if nextToken != nil {
		input = &ecs.ListServicesInput{
			Cluster:   cluster.ClusterName,
			NextToken: nextToken,
		}
	}

	req := client.EcsClient.ListServicesRequest(input)
	result, err := req.Send(context.Background())

	if err != nil {
		return nil, err
	}

	outputs = append(outputs, result.ListServicesOutput)

	if result.NextToken != nil {
		return client.listEcsServices(cluster, result.NextToken, outputs)
	}
	return outputs, nil
}

func (client Client) describeEcsServices(cluster *ecs.Cluster, services []string) (*ecs.DescribeServicesOutput, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  cluster.ClusterName,
		Services: services,
	}

	req := client.EcsClient.DescribeServicesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return result.DescribeServicesOutput, nil
}
