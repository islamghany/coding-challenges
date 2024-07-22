package algorithms

import (
	"cmp"
)

func min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func pickPivot[T cmp.Ordered](arr []T, high, low int) int {
	mid := low + (high-low)/2
	a, b, c := arr[low], arr[mid], arr[high]
	// return the median of a, b, c
	if a < b {
		if b < c {
			return mid
		} else if a < c {
			return high
		} else {
			return low
		}
	} else {
		if a < c {
			return low
		} else if b < c {
			return high
		} else {
			return mid
		}
	}
}

func partition[T cmp.Ordered](arr []T, low, high int) {
	// if there are no elements or only one element in the array then return
	if low >= high {
		return
	}
	// pick a pivot element
	pivIdx := pickPivot(arr, high, low)
	piv := arr[pivIdx]
	// swap the pivot element with the last element
	arr[pivIdx], arr[high] = arr[high], arr[pivIdx]
	i := low
	for j := low; j < high; j++ {
		if arr[j] <= piv {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	// swap the pivot element with the element at index i
	arr[i], arr[high] = arr[high], arr[i]
	// recursively sort the two partitions
	partition(arr, low, i-1)
	partition(arr, i+1, high)

}

func QuickSort[T cmp.Ordered](arr []T) {
	partition(arr, 0, len(arr)-1)
}
