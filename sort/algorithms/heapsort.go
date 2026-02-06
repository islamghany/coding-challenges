package algorithms

import (
	"cmp"
)

func leftChild(i int) int {
	return 2*i + 1
}

func rightChild(i int) int {
	return 2*i + 2
}

func parent(i int) int {
	return (i - 1) / 2
}

func sink[T cmp.Ordered](arr []T, i int) {
	n := len(arr)
	for {
		leftIdx := leftChild(i)
		rightIdx := rightChild(i)
		idx := i
		if leftIdx < n && arr[leftIdx] < arr[idx] {
			idx = leftIdx
		}
		if rightIdx < n && arr[rightIdx] < arr[idx] {
			idx = rightIdx
		}
		if idx == i {
			break
		}
		arr[i], arr[idx] = arr[idx], arr[i]
		i = idx
	}
}

func float[T cmp.Ordered](arr []T, i int) {
	for i > 0 {
		p := parent(i)
		if arr[p] < arr[i] {
			break
		}
		arr[p], arr[i] = arr[i], arr[p]
		i = p
	}
}

func minHeapify[T cmp.Ordered](arr []T) {
	n := len(arr)
	for i := n/2 - 1; i >= 0; i-- {
		sink(arr, i)
	}
}

func insert[T cmp.Ordered](arr []T, val T) []T {
	arr = append(arr, val)
	float(arr, len(arr)-1)
	return arr
}

func removeTop[T cmp.Ordered](arr []T) T {
	n := len(arr)
	if n == 1 {
		return arr[0]
	}
	top := arr[0]
	arr[0] = arr[n-1]
	arr = arr[:n-1]
	sink(arr, 0)
	return top
}

func Heapsort[T cmp.Ordered](arr []T) {
	n := len(arr)
	a := make([]T, n)
	minHeapify(arr)
	for i := 0; i < n; i++ {
		t := removeTop(arr)
		a[i] = t
	}
	copy(arr, a)
}
