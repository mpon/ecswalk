package awsapi

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
)

// EcsTaskInfo represents ...
type EcsTaskInfo struct {
	Task              ecs.Task
	ContainerInstance *ecs.ContainerInstance
	Instance          *ec2.Instance
}

// ShortTaskArn returns short task arn
func (t EcsTaskInfo) ShortTaskArn() string {
	return ShortArn(*t.Task.TaskArn)
}

// ShortTaskDefinitionArn return short task arn
func (t EcsTaskInfo) ShortTaskDefinitionArn() string {
	return ShortArn(*t.Task.TaskDefinitionArn)
}

// NewEcsTaskInfoList ...
func NewEcsTaskInfoList(tasks []ecs.Task) EcsTaskInfoList {
	var taskInfoList EcsTaskInfoList
	for _, t := range tasks {
		taskInfoList = append(taskInfoList, &EcsTaskInfo{
			Task: t,
		})
	}
	return taskInfoList
}

// ContainerInstanceArns returns distinct arns
func (tList EcsTaskInfoList) ContainerInstanceArns() []string {
	var arns []string
	for _, t := range tList {
		arns = append(arns, *t.Task.ContainerInstanceArn)
	}

	return sliceutil.DistinctSlice(arns)
}

// Ec2InstanceIds returns distinct instance ids
func (tList EcsTaskInfoList) Ec2InstanceIds() []string {
	var ids []string
	for _, t := range tList {
		ids = append(ids, *t.ContainerInstance.Ec2InstanceId)
	}

	return sliceutil.DistinctSlice(ids)
}

// SetContainerInstances ...
func (tList EcsTaskInfoList) SetContainerInstances(containerInstances []ecs.ContainerInstance) {
	for _, t := range tList {
		for _, c := range containerInstances {
			if *c.ContainerInstanceArn == *t.Task.ContainerInstanceArn {
				t.ContainerInstance = &c
			}
		}
	}
}

// SetEc2Instances ...
func (tList EcsTaskInfoList) SetEc2Instances(instances []ec2.Instance) {
	for _, t := range tList {
		for _, i := range instances {
			if *i.InstanceId == *t.ContainerInstance.Ec2InstanceId {
				t.Instance = &i
			}
		}
	}
}

// EcsTaskInfoList slice
type EcsTaskInfoList []*EcsTaskInfo

func (tList EcsTaskInfoList) Len() int {
	return len(tList)
}

func (tList EcsTaskInfoList) Swap(i, j int) {
	tList[i], tList[j] = tList[j], tList[i]
}

func (tList EcsTaskInfoList) Less(i, j int) bool {
	return *tList[i].Task.TaskArn < *tList[j].Task.TaskArn
}
