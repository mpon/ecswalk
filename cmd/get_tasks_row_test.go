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
