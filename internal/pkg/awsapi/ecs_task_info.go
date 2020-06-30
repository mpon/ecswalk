package awsapi

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/mpon/ecswalk/internal/pkg/sliceutil"
)

// ECSTaskInfo represents ...
type ECSTaskInfo struct {
	Task              ecs.Task
	ContainerInstance *ecs.ContainerInstance
	Instance          *ec2.Instance
}

// ShortTaskArn returns short task arn
func (t ECSTaskInfo) ShortTaskArn() string {
	return ShortArn(*t.Task.TaskArn)
}

// ShortTaskDefinitionArn return short task arn
func (t ECSTaskInfo) ShortTaskDefinitionArn() string {
	return ShortArn(*t.Task.TaskDefinitionArn)
}

// NewECSTaskInfoList ...
func NewECSTaskInfoList(tasks []ecs.Task) ECSTaskInfoList {
	var taskInfoList ECSTaskInfoList
	for _, t := range tasks {
		taskInfoList = append(taskInfoList, &ECSTaskInfo{
			Task: t,
		})
	}
	return taskInfoList
}

// ContainerInstanceArns returns distinct arns
func (tList ECSTaskInfoList) ContainerInstanceArns() []string {
	var arns []string
	for _, t := range tList {
		arns = append(arns, *t.Task.ContainerInstanceArn)
	}

	return sliceutil.DistinctSlice(arns)
}

// EC2InstanceIds returns distinct instance ids
func (tList ECSTaskInfoList) EC2InstanceIds() []string {
	var ids []string
	for _, t := range tList {
		ids = append(ids, *t.ContainerInstance.Ec2InstanceId)
	}

	return sliceutil.DistinctSlice(ids)
}

// SetContainerInstances ...
func (tList ECSTaskInfoList) SetContainerInstances(containerInstances []ecs.ContainerInstance) {
	for _, t := range tList {
		for _, c := range containerInstances {
			if *c.ContainerInstanceArn == *t.Task.ContainerInstanceArn {
				t.ContainerInstance = &c
			}
		}
	}
}

// SetEC2Instances ...
func (tList ECSTaskInfoList) SetEC2Instances(instances []ec2.Instance) {
	for _, t := range tList {
		for _, i := range instances {
			if *i.InstanceId == *t.ContainerInstance.Ec2InstanceId {
				t.Instance = &i
			}
		}
	}
}

// ECSTaskInfoList slice
type ECSTaskInfoList []*ECSTaskInfo

func (tList ECSTaskInfoList) Len() int {
	return len(tList)
}

func (tList ECSTaskInfoList) Swap(i, j int) {
	tList[i], tList[j] = tList[j], tList[i]
}

func (tList ECSTaskInfoList) Less(i, j int) bool {
	return *tList[i].Task.TaskArn < *tList[j].Task.TaskArn
}
