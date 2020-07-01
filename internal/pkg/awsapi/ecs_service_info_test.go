package awsapi

import (
	"fmt"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func TestEcsServiceInfoRreturnsNewEcsServiceInfoList(t *testing.T) {
	services := []ecs.Service{
		{
			ServiceName:    aws.String("name999"),
			TaskDefinition: aws.String("arn:aws:ecs:ap-northeast-1:123456789012:task-definition/task1:100"),
		},
		{
			ServiceName:    aws.String("name000"),
			TaskDefinition: aws.String("arn:aws:ecs:ap-northeast-1:123456789012:task-definition/task2:200"),
		},
	}
	taskDefinitions := []ecs.TaskDefinition{
		{
			Family:            aws.String("task2"),
			TaskDefinitionArn: aws.String("arn:aws:ecs:ap-northeast-1:123456789012:task-definition/task2:200"),
		},
		{
			Family:            aws.String("task1"),
			TaskDefinitionArn: aws.String("arn:aws:ecs:ap-northeast-1:123456789012:task-definition/task1:100"),
		},
	}

	sut := NewEcsServiceInfoList(services, taskDefinitions)

	if len(sut) != 2 {
		t.Fatalf("Failed New List %d", len(sut))
	}

	if *sut[0].Service.ServiceName != "name000" {
		t.Fatalf("Failed New List %s", *sut[0].Service.ServiceName)
	}
	if *sut[0].TaskDefinition.Family != "task2" {
		t.Fatalf("Failed New List %s", *sut[0].TaskDefinition.Family)
	}
}

func TestEcsServiceInfoSortByNameAsc(t *testing.T) {

	info := func(index int64) EcsServiceInfo {
		return EcsServiceInfo{
			Service: ecs.Service{
				ServiceName: aws.String(fmt.Sprintf("A%d", index)),
			},
			TaskDefinition: ecs.TaskDefinition{},
		}
	}

	list := EcsServiceInfoList{}

	for i := 3; i > 0; i-- {
		list = append(list, info(int64(i)))
	}

	sort.Sort(list)

	for i, v := range list {
		expect := info(int64(i + 1))
		if *v.Service.ServiceName != *expect.Service.ServiceName {
			t.Fatalf("Not orderd Name by asc\nexpect %#v\nbut %#v", expect, v)
		}
	}

}

func TestEcsServiceInfoReturnsTaskDefinitionArn(t *testing.T) {
	info := EcsServiceInfo{
		Service: ecs.Service{},
		TaskDefinition: ecs.TaskDefinition{
			TaskDefinitionArn: aws.String("arn:aws:ecs:ap-northeast-1:123456789012:task-definition/image:999"),
		},
	}

	sut := info.TaskDefinitionArn()
	expect := "image:999"

	if sut != expect {
		t.Fatalf("TaskDefinitionArn does not return correct format expect %s, but %s", expect, sut)
	}
}

func TestEcsServiceInfoRreturnsDockerImageInfo(t *testing.T) {
	info := EcsServiceInfo{
		Service: ecs.Service{},
		TaskDefinition: ecs.TaskDefinition{
			ContainerDefinitions: []ecs.ContainerDefinition{
				{
					Image: aws.String("example.com/image:abcd1234"),
				},
			},
		},
	}

	img := info.DockerImageName()
	expectImg := "image"
	tag := info.DockerImageTag()
	expectTag := "abcd1234"

	if img != expectImg {
		t.Fatalf("DockerImageName does not return correct format expect %s but %s", expectImg, img)
	}

	if tag != expectTag {
		t.Fatalf("DockerImageTag does not return correct format expect %s but %s", expectTag, tag)
	}
}

func TestEcsServiceInfoRreturnsMultipleDockerImageInfo(t *testing.T) {
	info := EcsServiceInfo{
		Service: ecs.Service{},
		TaskDefinition: ecs.TaskDefinition{
			ContainerDefinitions: []ecs.ContainerDefinition{
				{
					Image: aws.String("example.com/image1:abcd1234"),
				},
				{
					Image: aws.String("example.com/image2:1234abcd"),
				},
			},
		},
	}

	img := info.DockerImageName()
	expectImg := "image1,image2"
	tag := info.DockerImageTag()
	expectTag := "abcd1234,1234abcd"

	if img != expectImg {
		t.Fatalf("DockerImageName does not return correct format expect %s but %s", expectImg, img)
	}

	if tag != expectTag {
		t.Fatalf("DockerImageTag does not return correct format expect %s but %s", expectTag, tag)
	}
}
