package awsapi

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
)

// ECSContainerInstanceInfo represents ...
type ECSContainerInstanceInfo struct {
	ContainerInstance ecs.ContainerInstance
	Ec2Instance       ec2.Instance
}

// MemoryAvailable returns remaining cpu available
func (info ECSContainerInstanceInfo) MemoryAvailable() *int64 {
	for _, r := range info.ContainerInstance.RemainingResources {
		if *r.Name == "MEMORY" {
			return r.IntegerValue
		}
	}
	return nil
}

// CPUAvailable returns remaining cpu available
func (info ECSContainerInstanceInfo) CPUAvailable() *int64 {
	for _, r := range info.ContainerInstance.RemainingResources {
		if *r.Name == "CPU" {
			return r.IntegerValue
		}
	}
	return nil
}

// ShortContainerInstanceArn return short container instance arn
func (info ECSContainerInstanceInfo) ShortContainerInstanceArn() string {
	return ShortArn(*info.ContainerInstance.ContainerInstanceArn)
}

// DockerVersion returns docker version
func (info ECSContainerInstanceInfo) DockerVersion() string {
	return strings.Replace(*info.ContainerInstance.VersionInfo.DockerVersion, "DockerVersion: ", "", -1)
}

// ECSContainerInstanceInfoList slice
type ECSContainerInstanceInfoList []*ECSContainerInstanceInfo

// NewECSContainerInstanceInfoList constructor
func NewECSContainerInstanceInfoList(containerInstances []ecs.ContainerInstance) ECSContainerInstanceInfoList {
	var list ECSContainerInstanceInfoList
	for _, c := range containerInstances {
		list = append(list, &ECSContainerInstanceInfo{
			ContainerInstance: c,
		})
	}
	return list
}

// Ec2InstanceIds returns
func (cList ECSContainerInstanceInfoList) Ec2InstanceIds() []string {
	var ids []string
	for _, c := range cList {
		ids = append(ids, *c.ContainerInstance.Ec2InstanceId)
	}
	return sliceutil.DistinctSlice(ids)
}

// SetEC2Instances ...
func (cList ECSContainerInstanceInfoList) SetEC2Instances(ec2Instances []ec2.Instance) {
	for _, c := range cList {
		for _, i := range ec2Instances {
			if *i.InstanceId == *c.ContainerInstance.Ec2InstanceId {
				c.Ec2Instance = i
			}
		}
	}
}
