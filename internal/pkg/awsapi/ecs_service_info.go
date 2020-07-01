package awsapi

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// EcsServiceInfo represents ...
type EcsServiceInfo struct {
	Service        ecs.Service
	TaskDefinition ecs.TaskDefinition
}

// TaskDefinitionArn return ECS task definition short ARN
func (s EcsServiceInfo) TaskDefinitionArn() string {
	return ShortArn(*s.TaskDefinition.TaskDefinitionArn)
}

// DockerImageName return docker image name
func (s EcsServiceInfo) DockerImageName() string {
	var names []string
	for _, c := range s.TaskDefinition.ContainerDefinitions {
		image, _ := ShortDockerImage(*c.Image)
		names = append(names, image)
	}
	return strings.Join(names, ",")
}

// DockerImageTag return docker image tag
func (s EcsServiceInfo) DockerImageTag() string {
	var tags []string
	for _, c := range s.TaskDefinition.ContainerDefinitions {
		_, tag := ShortDockerImage(*c.Image)
		tags = append(tags, tag)
	}
	return strings.Join(tags, ",")
}

// EcsServiceInfoList slice
type EcsServiceInfoList []EcsServiceInfo

func (s EcsServiceInfoList) Len() int {
	return len(s)
}

func (s EcsServiceInfoList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s EcsServiceInfoList) Less(i, j int) bool {
	return *s[i].Service.ServiceName < *s[j].Service.ServiceName
}

// NewEcsServiceInfoList create ECS service infomation list
func NewEcsServiceInfoList(services []ecs.Service, taskDefinitions []ecs.TaskDefinition) EcsServiceInfoList {
	list := EcsServiceInfoList{}

	for _, s := range services {
		td := findTaskDefinition(s, taskDefinitions)
		list = append(list, EcsServiceInfo{
			Service:        s,
			TaskDefinition: *td,
		})
	}
	sort.Sort(list)
	return list
}

func findTaskDefinition(service ecs.Service, taskDefinitions []ecs.TaskDefinition) *ecs.TaskDefinition {
	for _, t := range taskDefinitions {
		if *service.TaskDefinition == *t.TaskDefinitionArn {
			return &t
		}
	}
	return nil
}
