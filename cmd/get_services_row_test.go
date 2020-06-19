package cmd

import (
	"fmt"
	"sort"
	"testing"
)

func TestGetServiceRowSortByNameAsc(t *testing.T) {

	row := func(index int64) GetServiceRow {
		return GetServiceRow{
			Name:           fmt.Sprintf("A%d", index),
			TaskDefinition: fmt.Sprintf("B%d", index),
			Image:          fmt.Sprintf("C%d", index),
			Tag:            fmt.Sprintf("D%d", index),
			DesiredCount:   2 * index,
			RunningCount:   1 * index,
		}
	}

	rows := GetServiceRows{}

	for i := 3; i > 0; i-- {
		rows = append(rows, row(int64(i)))
	}

	for i, v := range rows {
		expect := row(int64(3 - i))
		if v != expect {
			t.Fatalf("expect %#v\nbut %#v", expect, v)
		}
	}

	sort.Sort(rows)

	for i, v := range rows {
		expect := row(int64(i + 1))
		if v != expect {
			t.Fatalf("Not orderd Name by asc\nexpect %#v\nbut %#v", expect, v)
		}
	}

}
