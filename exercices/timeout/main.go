package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
* Poderia implementar usando context, mas o exercicio pede o uso do select,
mas eu ainda acho que é possível. Vou tentar.
Fiz a implementação em ./timeout_context
*/
func main() {
	start := time.Now()
	ch := make(chan string)

	go callAPISlow(ch)

	select {
	case r := <-ch:
		fmt.Println("r=", r)
		fmt.Println("took=", time.Since(start))
	case <-time.After(time.Second * 2):
		fmt.Println("ERRO: time exceeds 2 seconds.")
		fmt.Println("took=", time.Since(start))
		return
	}

}

func callAPISlow(ch chan<- string) {

	delay := rand.Intn(3)

	time.Sleep(time.Second * time.Duration(delay))

	ch <- "Response API"
}
