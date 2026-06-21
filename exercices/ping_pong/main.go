package main

import "fmt"

func main() {

	ch := make(chan string)
	done := make(chan bool)

	go ping(ch, done)
	go pong(ch)

	<-done
	fmt.Println("game done")

}

func ping(ch chan string, done chan bool) {
	for i := 1; i <= 5; i++ {
		ping := fmt.Sprintf("ping %d", i)
		fmt.Printf("[PING] sending %s\n", ping)

		ch <- ping

		pong := <-ch
		fmt.Printf("[PING] received: %s\n", pong)

	}
	close(ch)
	done <- true

}
func pong(ch chan string) {
	i := 1
	for m := range ch {
		pong := fmt.Sprintf("pong %d", i)
		fmt.Printf("[PONG] received: %s\n", m)
		fmt.Printf("[PONG] sending: %s\n", pong)
		ch <- pong
		i++
	}
}
