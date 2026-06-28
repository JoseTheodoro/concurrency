package main

import (
	"fmt"
)

func main() {
	var total int
	ch := make(chan int)
	fullSlice := createFullSlice()
	partOne := fullSlice[0:50]
	partTwo := fullSlice[50:]

	go sumSlice(partOne, ch)
	go sumSlice(partTwo, ch)

	for i := 1; i <= 2; i++ {
		total += <-ch
	}

	close(ch)

	fmt.Printf("Resultado da soma concorrente: %d\n", total)
}

func sumSlice(s []int, ch chan<- int) {
	var sum int
	for _, v := range s {
		sum += v
	}
	ch <- sum
}

func createFullSlice() []int {
	// criar um slice de 100 números inteiros
	var s []int

	for i := 1; i <= 100; i++ {
		s = append(s, i)
	}

	return s

}
