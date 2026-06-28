package main

import (
	"fmt"
	"time"
)

const (
	MAXNUMBERS = 200
)

func main() {
	chNumbers := make(chan int)
	chDoubleNumbers := make(chan int)
	chDone := make(chan bool)

	go pub(chNumbers)
	go double(chNumbers, chDoubleNumbers)
	go sub(chDoubleNumbers, chDone)

	<-chDone
	fmt.Println("all jobs done")

}

func pub(ch chan<- int) {
	for i := 1; i <= MAXNUMBERS; i++ {
		ch <- i
	}
	close(ch)
}
func double(chRead <-chan int, chDouble chan<- int) {
	for number := range chRead {
		d := number * number
		chDouble <- d
	}
	close(chDouble)
}
func sub(ch <-chan int, chDone chan<- bool) {
	for n := range ch {
		time.Sleep(time.Millisecond * 100) // só pra simular um consumo devagar.
		fmt.Printf("Double= %d\n", n)
	}
	chDone <- true

}
