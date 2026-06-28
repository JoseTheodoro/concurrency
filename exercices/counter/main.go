package main

import (
	"fmt"
	"sync"
)

var counter int
var wg sync.WaitGroup
var mu sync.Mutex

func main() {

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go increment()
	}

	wg.Wait()

	fmt.Printf("Todos os jobs finalizaram e counter= %d\n", counter)

}

func increment() {
	mu.Lock()
	defer mu.Unlock()
	counter++
	wg.Done()
}
