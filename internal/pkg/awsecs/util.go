package awsecs

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// ShortArn returns last splited part by slash
func ShortArn(arn string) string {
	const slash = "/"
	parts := strings.Split(arn, slash)
	return parts[len(parts)-1]
}

// ShortDockerImage returns docker image name and tag
func ShortDockerImage(image string) (string, string) {
	const slash = "/"
	const colon = ":"
	parts := strings.Split(image, slash)
	names := strings.Split(parts[len(parts)-1], colon)

	if len(names) == 1 {
		return names[0], ""
	}
	return names[0], names[1]
}

// FindService to find service by task definition
func FindService(services []ecs.Service, taskDefinition string) ecs.Service {
	for _, service := range services {
		if *service.TaskDefinition == taskDefinition {
			return service
		}
	}
	return ecs.Service{}
}
