package algorithms

import "cmp"

func mergeSort[T cmp.Ordered](arr []T, low, high int) []T {
	if low >= high {
		return arr
	}
	mid := low + (high-low)/2
	left := mergeSort(arr, low, mid)
	right := mergeSort(arr, mid+1, high)
	i, j, k := low, mid+1, 0
	temp := make([]T, len(arr))
	for i <= mid && j <= high {
		if left[i] < right[j] {
			temp[k] = left[i]
			i++
		} else {
			temp[k] = right[j]
			j++
		}
		k++
	}
	for i <= mid {
		temp[k] = left[i]
		i++
		k++
	}
	for j <= high {
		temp[k] = right[j]
		j++
		k++
	}
	for i := 0; i < k; i++ {
		arr[low+i] = temp[i]
	}
	return arr
}

func MergeSort[T cmp.Ordered](arr []T) {
	// a := make([]T, len(arr))
	// copy(a, arr)
	mergeSort(arr, 0, len(arr)-1)

}
