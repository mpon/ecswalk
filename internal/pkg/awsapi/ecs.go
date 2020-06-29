package awsapi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

// GetECSCluster to get an ECS cluster
func (client Client) GetECSCluster(clusterName string) (*ecs.Cluster, error) {
	output, err := client.describeECSClusters()
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

// GetECSService to get an ECS Service
func (client Client) GetECSService(cluster *ecs.Cluster, serviceName string) (*ecs.Service, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  cluster.ClusterName,
		Services: []string{serviceName},
	}

	req := client.ECSClient.DescribeServicesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return &result.DescribeServicesOutput.Services[0], nil
}

// GetAllECSServices to get all ECS Services
func (client Client) GetAllECSServices(cluster *ecs.Cluster) ([]ecs.Service, error) {
	outputs, err := client.describeAllECSServices(cluster)
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

// GetAllECSClusters to get all ECS Clusters
func (client Client) GetAllECSClusters() ([]ecs.Cluster, error) {
	output, err := client.describeECSClusters()
	if err != nil {
		return nil, err
	}
	return output.Clusters, nil
}

// GetECSTaskDefinitions to get ECS task definition list
func (client Client) GetECSTaskDefinitions(cluster *ecs.Cluster, services []ecs.Service) ([]ecs.TaskDefinition, error) {
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

func (client Client) describeECSClusters() (*ecs.DescribeClustersOutput, error) {
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

func (client Client) describeAllECSServices(cluster *ecs.Cluster) ([]*ecs.DescribeServicesOutput, error) {
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

func (client Client) describeTaskDefinition(taskDefinitionArn string) (*ecs.DescribeTaskDefinitionOutput, error) {
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

// GetECSTasks to get ECS tasks of specified cluster and service
func (client Client) GetECSTasks(cluster *ecs.Cluster, service *ecs.Service) ([]ecs.Task, error) {
	input := &ecs.ListTasksInput{
		Cluster:     cluster.ClusterName,
		ServiceName: service.ServiceName,
	}

	req := client.ECSClient.ListTasksRequest(input)
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

	req := client.ECSClient.DescribeTasksRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, xerrors.Errorf("ECS DescribeTasks: %w", err)
	}
	return result.DescribeTasksOutput, nil
}

// ListAllContainerInstances to list all container instances
func (client Client) ListAllContainerInstances(clusetr *ecs.Cluster) (*ecs.ListContainerInstancesOutput, error) {
	input := &ecs.ListContainerInstancesInput{
		Cluster: clusetr.ClusterName,
	}

	req := client.ECSClient.ListContainerInstancesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return result.ListContainerInstancesOutput, nil
}

// DescribeContainerInstances to describe container instances
func (client Client) DescribeContainerInstances(cluster *ecs.Cluster, containerInstanceArns []string) (*ecs.DescribeContainerInstancesOutput, error) {
	input := &ecs.DescribeContainerInstancesInput{
		Cluster:            cluster.ClusterName,
		ContainerInstances: containerInstanceArns,
	}

	req := client.ECSClient.DescribeContainerInstancesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, xerrors.Errorf("ECS DescribeContainerInstances: %w", err)
	}
	return result.DescribeContainerInstancesOutput, nil
}

func (client Client) listECSServicesRecursively(cluster *ecs.Cluster) ([]*ecs.ListServicesOutput, error) {
	outputs, err := client.listECSServices(cluster, nil, nil)

	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func (client Client) listECSServices(cluster *ecs.Cluster, nextToken *string, outputs []*ecs.ListServicesOutput) ([]*ecs.ListServicesOutput, error) {
	input := &ecs.ListServicesInput{
		Cluster: cluster.ClusterName,
	}

	if nextToken != nil {
		input = &ecs.ListServicesInput{
			Cluster:   cluster.ClusterName,
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

func (client Client) describeECSServices(cluster *ecs.Cluster, services []string) (*ecs.DescribeServicesOutput, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  cluster.ClusterName,
		Services: services,
	}

	req := client.ECSClient.DescribeServicesRequest(input)
	result, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}
	return result.DescribeServicesOutput, nil
}
