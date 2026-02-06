package main

import (
	"fmt"
	"slices"
	"testing"
	"unixsort/algorithms"
)

func TestSortAlgorithms(t *testing.T) {
	testCases := []struct {
		name string
		arr  []int
	}{
		{
			name: "small array",
			arr:  []int{3, 2, 1, 5, 4, 6},
		},
		{
			name: "empty array",
			arr:  []int{},
		},
		{
			name: "sorted array",
			arr:  []int{1, 2, 3, 4, 5},
		},
		{
			name: "reverse sorted array",
			arr:  []int{5, 4, 3, 2, 1},
		},
		{
			name: "large array",
			arr:  []int{3, 2, 1, 5, 4, 6, 9, 8, 7, 10},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			arr := make([]int, len(tc.arr))
			copy(arr, tc.arr)

			algorithms.QuickSort(arr)
			if !slices.IsSorted(arr) {
				t.Errorf("QuickSort failed+%v", arr)
			}

			arr = make([]int, len(tc.arr))
			copy(arr, tc.arr)
			algorithms.MergeSort(arr)
			if !slices.IsSorted(arr) {
				t.Errorf("MergeSort failed%+v", arr)
			}

			arr = make([]int, len(tc.arr))
			copy(arr, tc.arr)
			algorithms.Heapsort(arr)
			if !slices.IsSorted(arr) {
				t.Errorf("HeapSort failed%+v", arr)
			}

			// arr = make([]int, len(tc.arr))
			// copy(arr, tc.arr)
			// algorithms.HeapSort(arr)
			// if !slices.IsSorted(arr) {
			// 	t.Errorf("RadixSort failed")
			// }
		})
	}
}

func TestBogoSort(t *testing.T) {
	arr := []int{3, 2, -1, 100, 122, 222, 4, 10, 5, 1, 4}
	fmt.Println("Unsorted array:", arr)
	algorithms.RandomSort(arr)
	fmt.Println("Sorted array:", arr)
	if !slices.IsSorted(arr) {
		t.Errorf("RandomSort failed%+v", arr)
	}
}
