package fuzzyfinder

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mpon/ecswalk/internal/pkg/awsapi"
)

// FindCluster find fuzzily ecs.Cluster
func FindCluster(clusters []ecs.Cluster) (*ecs.Cluster, error) {
	idx, err := fuzzyfinder.Find(clusters,
		func(i int) string {
			return *clusters[i].ClusterName
		},
		fuzzyfinder.WithPromptString("Select Cluster:"),
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			cluster := clusters[i]
			return fmt.Sprintf("%s\n\nServices: %d\nRunning Tasks: %d\nPending Tasks: %d\nContainer instances: %d",
				*cluster.ClusterName,
				*cluster.ActiveServicesCount,
				*cluster.RunningTasksCount,
				*cluster.PendingTasksCount,
				*cluster.RegisteredContainerInstancesCount,
			)
		}),
	)

	if err != nil {
		return nil, err
	}
	return &clusters[idx], nil
}

// FindService find fuzzily ecs.Service
func FindService(services []ecs.Service) (*ecs.Service, error) {
	idx, err := fuzzyfinder.Find(services,
		func(i int) string {
			s := services[i]
			return *s.ServiceName
		},
		fuzzyfinder.WithPromptString("Select Service:"),
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			s := services[i]
			return fmt.Sprintf("%s\n\nTask Definition: %s\nDesired tasks: %d\nRunning tasks: %d",
				*s.ServiceName,
				awsapi.ShortArn(*s.TaskDefinition),
				*s.DesiredCount,
				*s.RunningCount,
			)
		}),
	)

	if err != nil {
		return nil, err
	}

	return &services[idx], nil
}
