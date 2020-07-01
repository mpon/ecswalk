package awsapi

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
)

// EcsContainerInstanceInfo represents ...
type EcsContainerInstanceInfo struct {
	ContainerInstance ecs.ContainerInstance
	Ec2Instance       ec2.Instance
}

// MemoryAvailable returns remaining cpu available
func (info EcsContainerInstanceInfo) MemoryAvailable() *int64 {
	for _, r := range info.ContainerInstance.RemainingResources {
		if *r.Name == "MEMORY" {
			return r.IntegerValue
		}
	}
	return nil
}

// CPUAvailable returns remaining cpu available
func (info EcsContainerInstanceInfo) CPUAvailable() *int64 {
	for _, r := range info.ContainerInstance.RemainingResources {
		if *r.Name == "CPU" {
			return r.IntegerValue
		}
	}
	return nil
}

// ShortContainerInstanceArn return short container instance arn
func (info EcsContainerInstanceInfo) ShortContainerInstanceArn() string {
	return ShortArn(*info.ContainerInstance.ContainerInstanceArn)
}

// DockerVersion returns docker version
func (info EcsContainerInstanceInfo) DockerVersion() string {
	return strings.Replace(*info.ContainerInstance.VersionInfo.DockerVersion, "DockerVersion: ", "", -1)
}

// EcsContainerInstanceInfoList slice
type EcsContainerInstanceInfoList []*EcsContainerInstanceInfo

// NewEcsContainerInstanceInfoList constructor
func NewEcsContainerInstanceInfoList(containerInstances []ecs.ContainerInstance) EcsContainerInstanceInfoList {
	var list EcsContainerInstanceInfoList
	for _, c := range containerInstances {
		list = append(list, &EcsContainerInstanceInfo{
			ContainerInstance: c,
		})
	}
	return list
}

// Ec2InstanceIds returns
func (cList EcsContainerInstanceInfoList) Ec2InstanceIds() []string {
	var ids []string
	for _, c := range cList {
		ids = append(ids, *c.ContainerInstance.Ec2InstanceId)
	}
	return sliceutil.DistinctSlice(ids)
}

// SetEc2Instances ...
func (cList EcsContainerInstanceInfoList) SetEc2Instances(ec2Instances []ec2.Instance) {
	for _, c := range cList {
		for _, i := range ec2Instances {
			if *i.InstanceId == *c.ContainerInstance.Ec2InstanceId {
				c.Ec2Instance = i
			}
		}
	}
}
