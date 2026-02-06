package algorithms

import (
	"cmp"
	"math/rand"
	"time"
)

// isSorted checks if the array is sorted
func isSorted[T cmp.Ordered](arr []T) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i-1] > arr[i] {
			return false
		}
	}
	return true
}

// shuffle randomly shuffles the array
func shuffle[T cmp.Ordered](arr []T) {
	rand.Seed(time.Now().UnixNano())
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// bogoSort sorts the array using BogoSort
func bogoSort[T cmp.Ordered](arr []T) {
	for !isSorted(arr) {
		shuffle(arr)
	}
}

func RandomSort[T cmp.Ordered](arr []T) {
	bogoSort(arr)
}
