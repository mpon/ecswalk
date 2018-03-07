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
