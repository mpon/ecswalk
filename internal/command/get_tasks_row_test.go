package command

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestGetTaskRowSortByTaskIdAsc(t *testing.T) {

	row := func(index int64) GetTaskRow {
		return GetTaskRow{
			TaskID:         fmt.Sprintf("A%d", index),
			TaskDefinition: fmt.Sprintf("B%d", index),
			Status:         fmt.Sprintf("C%d", index),
			PrivateIP:      fmt.Sprintf("D%d", index),
		}
	}

	rows := GetTaskRows{}

	for i := 3; i > 0; i-- {
		r := row(int64(i))
		rows = append(rows, &r)
	}

	for i, v := range rows {
		r := row(int64(3 - i))
		expect := &r
		if !reflect.DeepEqual(v, expect) {
			t.Fatalf("expect %#v\nbut %#v", expect, v)
		}
	}

	sort.Sort(rows)

	for i, v := range rows {
		r := row(int64(i + 1))
		expect := &r
		if !reflect.DeepEqual(v, expect) {
			t.Fatalf("Not orderd TaskID by asc\nexpect %#v\nbut %#v", expect, v)
		}
	}

}
