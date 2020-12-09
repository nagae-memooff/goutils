package utils

import (
	"fmt"
	"testing"
)

func TestSplitToIntSlice(t *testing.T) {

	CheckSlice(t, "1,2,3", []int{1, 2, 3})
	CheckSlice(t, "", []int{})
	CheckSlice(t, " ", []int{})
	CheckSlice(t, "1,2,", []int{1, 2})
	CheckSlice(t, " 1,2, ", []int{1, 2})
	CheckSlice(t, " 1,2,3 ", []int{1, 2, 3})
	CheckSlice(t, " 1,2,3 ,", []int{1, 2, 3})
	CheckSlice(t, " 1,2,3, ", []int{1, 2, 3})
	CheckSlice(t, " 1,2,3, 4, ", []int{1, 2, 3, 4})
}

func CheckSlice(t *testing.T, str string, expect_slice []int) bool {
	slice := SplitToIntSlice(str)

	slice_string := fmt.Sprintf("%v", slice)
	expect_slice_string := fmt.Sprintf("%v", expect_slice)

	fmt.Printf("str: '%s' slice: %v(%d), expect_slice: %v(%d).\n", str, slice, len(slice), expect_slice, len(expect_slice))

	if slice_string == expect_slice_string {
		return true
	} else {
		t.Errorf("x: %v, expect: %v.\n", slice, expect_slice)
		return false
	}
}
