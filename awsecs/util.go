// Copyright © 2018 Masato Oshima
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
