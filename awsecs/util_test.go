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
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func TestShortArn(t *testing.T) {
	arn := "arn:aws:ecs:region:account-id:cluster/cluster-name"
	expect := "cluster-name"
	result := ShortArn(arn)

	if result != expect {
		t.Fatalf("expect %s\nbut %s", expect, result)
	}
}

func TestShortArnInvalidArguments(t *testing.T) {
	arn := "invalid"
	expect := "invalid"
	result := ShortArn(arn)

	if result != expect {
		t.Fatalf("expect %s\nbut %s", expect, result)
	}
}

func TestShortDockerImage(t *testing.T) {
	image := "us.gcr.io/my-project/my-image:test"
	expectImage := "my-image"
	expectTag := "test"
	resultImage, resultTag := ShortDockerImage(image)

	if resultImage != expectImage {
		t.Fatalf("expect %s\nbut %s", expectImage, resultImage)
	}

	if resultTag != expectTag {
		t.Fatalf("expect %s\nbut %s", expectTag, resultTag)
	}
}

func TestShortDockerImageWithoutTag(t *testing.T) {
	image := "us.gcr.io/my-project/my-image"
	expectImage := "my-image"
	expectTag := ""
	resultImage, resultTag := ShortDockerImage(image)

	if resultImage != expectImage {
		t.Fatalf("expect %s\nbut %s", expectImage, resultImage)
	}

	if resultTag != expectTag {
		t.Fatalf("expect %s\nbut %s", expectTag, resultTag)
	}
}

func TestFindService(t *testing.T) {
	taskDefinition := "arn:aws:ecs:us-east-1:123456789012:task-definition/hello_world:8"
	services := []ecs.Service{
		ecs.Service{
			TaskDefinition: aws.String(taskDefinition),
		},
		ecs.Service{
			TaskDefinition: aws.String("arn:aws:ecs:us-east-1:123456789012:task-definition/hello_world:9"),
		},
	}

	expect := services[0]
	result := FindService(services, taskDefinition)

	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("expect %#v\nbut %#v", expect, result)
	}
}

func TestFindServiceNothing(t *testing.T) {
	taskDefinition := "arn:aws:ecs:us-east-1:123456789012:task-definition/hello_world:8"
	services := []ecs.Service{
		ecs.Service{
			TaskDefinition: aws.String("arn:aws:ecs:us-east-1:123456789012:task-definition/hello_world:7"),
		},
		ecs.Service{
			TaskDefinition: aws.String("arn:aws:ecs:us-east-1:123456789012:task-definition/hello_world:9"),
		},
	}

	expect := ecs.Service{}
	result := FindService(services, taskDefinition)

	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("expect %#v\nbut %#v", expect, result)
	}
}
