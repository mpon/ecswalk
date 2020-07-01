package awsapi

import (
	"fmt"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func TestGetTaskRowSortByTaskIdAsc(t *testing.T) {

	info := func(index int64) EcsTaskInfo {
		return EcsTaskInfo{
			Task: ecs.Task{
				TaskArn:           aws.String(fmt.Sprintf("A%d", index)),
				TaskDefinitionArn: aws.String(fmt.Sprintf("B%d", index)),
			},
		}
	}

	tList := EcsTaskInfoList{}

	for i := 3; i > 0; i-- {
		x := info(int64(i))
		tList = append(tList, &x)
	}

	sort.Sort(tList)

	for i, v := range tList {
		x := info(int64(i + 1))
		expect := x
		if v.Task.TaskDefinitionArn == expect.Task.TaskDefinitionArn {
			t.Fatalf("Not orderd by asc\nexpect %#v\nbut %#v", expect, v)
		}
	}

}
