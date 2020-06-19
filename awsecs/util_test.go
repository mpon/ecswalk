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
		{
			TaskDefinition: aws.String(taskDefinition),
		},
		{
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
		{
			TaskDefinition: aws.String("arn:aws:ecs:us-east-1:123456789012:task-definition/hello_world:7"),
		},
		{
			TaskDefinition: aws.String("arn:aws:ecs:us-east-1:123456789012:task-definition/hello_world:9"),
		},
	}

	expect := ecs.Service{}
	result := FindService(services, taskDefinition)

	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("expect %#v\nbut %#v", expect, result)
	}
}
