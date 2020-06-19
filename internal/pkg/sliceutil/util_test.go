package sliceutil

import (
	"reflect"
	"testing"
)

func TestChunkedSlice(t *testing.T) {
	slice := []string{"a", "b", "c", "a", "b"}
	expect := [][]string{
		[]string{"a", "b", "c"},
		[]string{"a", "b"},
	}
	result := ChunkedSlice(slice, 3)

	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("expect %s\nbut %s", expect, result)
	}
}

func TestDistinctSlice(t *testing.T) {
	slice := []string{"a", "b", "c", "a", "a", "b", "d"}
	expect := []string{"a", "b", "c", "d"}
	result := DistinctSlice(slice)

	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("expect %s\nbut %s", expect, result)
	}
}
