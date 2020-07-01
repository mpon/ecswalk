package awsapi

import (
	"strings"
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
