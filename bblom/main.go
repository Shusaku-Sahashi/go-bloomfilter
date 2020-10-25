package main

import (
	"chace-alg/bblom/src"
	"fmt"
)

func main() {
	bloomFilter := src.NewBloom(1000, 0.01)

	bloomFilter.Add(200)
	fmt.Println(bloomFilter.Exist(200))
}
