package util_test

import (
	"testing"

	"github.com/Sheriff-Hoti/paper-tui/util"
)

func TestGet(t *testing.T) {
	nums := []int{10, 20, 30}

	// Valid index
	val, ok := util.Get(nums, 1)
	if !ok || val != 20 {
		t.Errorf("Get failed: expected (20,true), got (%v,%v)", val, ok)
	}

	// Negative index
	val, ok = util.Get(nums, -1)
	if ok || val != 0 {
		t.Errorf("Get failed for negative index: expected (0,false), got (%v,%v)", val, ok)
	}

	// Out-of-range index
	val, ok = util.Get(nums, 5)
	if ok || val != 0 {
		t.Errorf("Get failed for out-of-range index: expected (0,false), got (%v,%v)", val, ok)
	}

	// Nil slice
	var nilSlice []int
	val, ok = util.Get(nilSlice, 0)
	if ok || val != 0 {
		t.Errorf("Get failed for nil slice: expected (0,false), got (%v,%v)", val, ok)
	}

	// Slice of pointers
	ptrs := []*int{nil, new(int)}
	vp, ok := util.Get(ptrs, 0)
	if !ok || vp != nil {
		t.Errorf("Get failed for pointer slice: expected (nil,true), got (%v,%v)", vp, ok)
	}
}

func TestGetOrDefault(t *testing.T) {
	nums := []int{10, 20, 30}

	// Valid index
	val := util.GetOrDefault(nums, 2, 99)
	if val != 30 {
		t.Errorf("GetOrDefault failed: expected 30, got %v", val)
	}

	// Out-of-range index
	val = util.GetOrDefault(nums, 5, 99)
	if val != 99 {
		t.Errorf("GetOrDefault failed: expected 99, got %v", val)
	}

	// Negative index
	val = util.GetOrDefault(nums, -1, 42)
	if val != 42 {
		t.Errorf("GetOrDefault failed for negative index: expected 42, got %v", val)
	}

	// Nil slice
	var nilSlice []int
	val = util.GetOrDefault(nilSlice, 0, 7)
	if val != 7 {
		t.Errorf("GetOrDefault failed for nil slice: expected 7, got %v", val)
	}

	// Slice of pointers
	ptrs := []*int{nil, new(int)}
	vp := util.GetOrDefault(ptrs, 0, new(int))
	if vp != nil {
		t.Errorf("GetOrDefault failed for pointer slice: expected nil, got %v", vp)
	}
}

func TestSet(t *testing.T) {
	nums := []int{1, 2, 3}

	// Valid index
	ok := util.Set(nums, 1, 99)
	if !ok || nums[1] != 99 {
		t.Errorf("Set failed: expected nums[1]=99, got %v, ok=%v", nums[1], ok)
	}

	// Out-of-range index
	ok = util.Set(nums, 5, 42)
	if ok {
		t.Errorf("Set should fail for out-of-range index, but ok=%v", ok)
	}

	// Negative index
	ok = util.Set(nums, -1, 42)
	if ok {
		t.Errorf("Set should fail for negative index, but ok=%v", ok)
	}

	// Nil slice
	var nilSlice []int
	ok = util.Set(nilSlice, 0, 10)
	if ok {
		t.Errorf("Set should fail for nil slice, but ok=%v", ok)
	}

	// Slice of pointers
	ptrs := []*int{nil, new(int)}
	ok = util.Set(ptrs, 0, new(int))
	if !ok || ptrs[0] == nil {
		t.Errorf("Set failed for pointer slice, expected non-nil, got %v", ptrs[0])
	}
}
