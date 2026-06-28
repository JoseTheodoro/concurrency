package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

func main() {
	ch := make(chan string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	go callSomeAPI(ctx, ch)

	select {
	case r := <-ch:
		fmt.Println("API RESPONSE=", r)
		close(ch)
	case <-ctx.Done():
		fmt.Println("Error: API time execeeds 2 seconds.")
	}
}

func callSomeAPI(ctx context.Context, ch chan<- string) {

	delay := rand.IntN(4) + 1

	select {
	/*
		* se fosse uma chama real para http, NewRequestWithContext abstrai
		essa lógica do select.
	*/
	case <-time.After(time.Second * time.Duration(delay)):
		ch <- "Response API Done!"
	case <-ctx.Done():
		fmt.Println("Erro: timeout na API")
	}

}
