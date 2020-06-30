package awsapi

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// ECSServiceInfo represents ...
type ECSServiceInfo struct {
	Service        ecs.Service
	TaskDefinition ecs.TaskDefinition
}

// for _, containerDefinition := range td.ContainerDefinitions {
// 	image, tag := awsapi.ShortDockerImage(*containerDefinition.Image)
// 	rows = append(rows, GetServiceRow{
// 		Name:           *service.ServiceName,
// 		TaskDefinition: awsapi.ShortArn(*td.TaskDefinitionArn),
// 		Image:          image,
// 		Tag:            tag,
// 		DesiredCount:   *service.DesiredCount,
// 		RunningCount:   *service.RunningCount,
// 	})
// }

// TaskDefinitionArn return ECS task definition short ARN
func (s ECSServiceInfo) TaskDefinitionArn() string {
	return ShortArn(*s.TaskDefinition.TaskDefinitionArn)
}

// DockerImageName return docker image name
func (s ECSServiceInfo) DockerImageName() string {
	var names []string
	for _, c := range s.TaskDefinition.ContainerDefinitions {
		image, _ := ShortDockerImage(*c.Image)
		names = append(names, image)
	}
	return strings.Join(names, ",")
}

// DockerImageTag return docker image tag
func (s ECSServiceInfo) DockerImageTag() string {
	var tags []string
	for _, c := range s.TaskDefinition.ContainerDefinitions {
		_, tag := ShortDockerImage(*c.Image)
		tags = append(tags, tag)
	}
	return strings.Join(tags, ",")
}

// ECSServiceInfoList slice
type ECSServiceInfoList []ECSServiceInfo

func (s ECSServiceInfoList) Len() int {
	return len(s)
}

func (s ECSServiceInfoList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ECSServiceInfoList) Less(i, j int) bool {
	return *s[i].Service.ServiceName < *s[j].Service.ServiceName
}

// NewECSServiceInfoList create ECS service infomation list
func NewECSServiceInfoList(services []ecs.Service, taskDefinitions []ecs.TaskDefinition) ECSServiceInfoList {
	list := ECSServiceInfoList{}

	for _, s := range services {
		td := findTaskDefinition(s, taskDefinitions)
		list = append(list, ECSServiceInfo{
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
