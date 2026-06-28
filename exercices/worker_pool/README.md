## Worker Pool
Padrão de pool de trabalhadores

Objetivo: Aprender a gerenciar uma quantidade fixa de goroutines consumindo dinamicamente de uma única fonte de trabalho.


1. Criar um canal de tarefas com buffer para 10 números inteiros e preencha-o com números de 1 a 10. 
2. Inicie 3 goroutines atuando como "trabalhadoras". 
3. Cada trabalhadora deve ler continuamente desse canal, pausar por meio segundo simulando um cálculo, e imprimir uma mensagem como "Trabalhador 2 finalizou a tarefa 5". 
4. O programa deve encerrar quando a fila esvaziar.

### O que foi legal nesse exercício?
Orquestrar trabalhadores, goroutines em uma única fila de trabalho foi bem legal para fixar conhecimento de orquestração.

A dificuldade que tive foi decidir em como fazer a main saber ou esperar todos os trabalhadores terminarem.

A primeira implementação use WaitGroups
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	task := make(chan int, 10)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, task, &wg)
	}

	for i := 1; i <= 10; i++ {
		task <- i
	}

	close(task)

	wg.Wait()

	fmt.Println("Workers finalizados")

}

func worker(id int, task <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range task {
		time.Sleep(time.Second * 1)
		fmt.Printf("Worker %d finalizou a tarefa %d\n", id, t)
	}

}

```

A segunda implementação foi criar um canal de tarefas concluídas
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	workers := 3
	tarefas := 10
	task := make(chan int, tarefas)
	done := make(chan bool, workers)

	for i := 1; i <= workers; i++ {
		go worker(i, task, done)
	}

	for i := 1; i <= tarefas; i++ {
		task <- i
	}

	close(task)

	for i := 1; i <= 3; i++ {
		<-done
	}

	fmt.Println("Workers finalizados")

}

func worker(id int, task <-chan int, d chan<- bool) {
	for t := range task {
		time.Sleep(time.Second * 1)
		fmt.Printf("Worker %d finalizou a tarefa %d\n", id, t)
	}
	d <- true
}
```
Particularmente falando, achei essa solução mais elegante. E purista com o modelo CSP.

Implementação final: parametrizei o número de trabalhadores, número de tarefas e coloquei em quanto tempo levou para todos os trabalhadores terminarem.
```go
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
```

output:
```bash
Worker 3 finalizou a tarefa 993
Worker 8 finalizou a tarefa 996
Worker 7 finalizou a tarefa 994
Worker 4 finalizou a tarefa 990
Worker 10 finalizou a tarefa 989
Worker 9 finalizou a tarefa 998
Worker 15 finalizou a tarefa 997
Worker 5 finalizou a tarefa 999
Workers finalizados em 50.0552s
```