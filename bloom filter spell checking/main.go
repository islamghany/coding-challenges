package main

import (
	"fmt"
	bloomfilter "spellchecker/bloom_filter"
)

func main() {
	bf := bloomfilter.NewBloomFilter(10, 0.0001)
	bf.Load()
	bf.Add([]byte("hello"))
	bf.Add([]byte("world"))
	bf.Add([]byte("good"))
	bf.Add([]byte("morning"))
	bf.Add([]byte("evening"))
	bf.Add([]byte("night"))

	fmt.Println(bf.Contains([]byte("hello")))
	fmt.Println(bf.Contains([]byte("world")))
	fmt.Println(bf.Contains([]byte("goo22d")))
	fmt.Println(bf.Contains([]byte("morning")))
	fmt.Println(bf.Contains([]byte("even2ing")))
	fmt.Println(bf.Contains([]byte("night")))

	bf.Save()
}
