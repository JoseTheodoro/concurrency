package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	workers := 20
	tarefas := 1000
	task := make(chan int, tarefas)
	done := make(chan bool, workers)

	for i := 1; i <= workers; i++ {
		go worker(i, task, done)
	}

	for i := 1; i <= tarefas; i++ {
		task <- i
	}

	close(task)

	for i := 1; i <= workers; i++ {
		<-done
	}

	took := time.Since(start)
	fmt.Printf("Workers finalizados em %v\n", took)

}

func worker(id int, task <-chan int, d chan<- bool) {
	for t := range task {
		time.Sleep(time.Second * 1)
		fmt.Printf("Worker %d finalizou a tarefa %d\n", id, t)
	}
	d <- true
}
