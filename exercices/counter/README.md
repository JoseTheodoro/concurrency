### Counter com Mutex

Objetivo: Forçar um race condition intencional e resolvê-la isolando a memória.

1. Criar uma variável global contador iniciando em 0. 
2. Monte um loop que dispare 1000 goroutines. 
3. A única tarefa de cada goroutine é incrementar o contador em 1. 
4. Irei usar sync.WaitGroup para garantir que a função main só imprima o valor após todas as execuções terminarem. 
5. Rodar primeiro sem controle de acesso para observar race condition
6. Depois realizar o fix para o race condition

### O que foi legal nesse exercício?

Simular o race condition e não acontecer um panic ou null pointer exception.
Ao executar o código com race consition feito propositalmente
```go
var counter int
var wg sync.WaitGroup

func main() {

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go increment()
	}

	wg.Wait()

	fmt.Printf("Todos os jobs finalizaram e counter= %d\n", counter)

}

func increment() {
	counter++
	wg.Done()
}
```

A cada execução counter tem um valor:
```bash
❯ go run main.go
Todos os jobs finalizaram e counter= 968
❯ go run main.go
Todos os jobs finalizaram e counter= 944
❯ go run main.go
Todos os jobs finalizaram e counter= 958
```

O race condition é a sobrescrita silenciosa na memória e quando uma ou mais goroutines leem o valor 5 ao mesmo tempo por exemplo, ambas somam para 6. Podem ter sido 10 goroutines que execuram a operação no ao mesmo tempo, mas o contador só andou 1.

Nota: Verifique se ha race condition executando seu programa com a flag -race
```bash
go run -race main.go
```

Esse comando vai executar seu programa e te ajudar a entender se sua implementação tem race condition, indicando exatamente a leitura ou escrita ocorreram simultaneamente.

Resultado:
```bash
❯ go run -race main.go
==================
WARNING: DATA RACE
Read at 0x0001011d7650 by goroutine 8:
  main.increment()
      /Users/josetheodoro/workspace/learning/golang/concurrency/exercices/counter/main.go:25 +0x28

Previous write at 0x0001011d7650 by goroutine 7:
  main.increment()
      /Users/josetheodoro/workspace/learning/golang/concurrency/exercices/counter/main.go:25 +0x40

Goroutine 8 (running) created at:
  main.main()
      /Users/josetheodoro/workspace/learning/golang/concurrency/exercices/counter/main.go:15 +0x44

Goroutine 7 (running) created at:
  main.main()
      /Users/josetheodoro/workspace/learning/golang/concurrency/exercices/counter/main.go:15 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x0001011d7650 by goroutine 13:
  main.increment()
      /Users/josetheodoro/workspace/learning/golang/concurrency/exercices/counter/main.go:25 +0x40

Previous write at 0x0001011d7650 by goroutine 19:
  main.increment()
      /Users/josetheodoro/workspace/learning/golang/concurrency/exercices/counter/main.go:25 +0x40

Goroutine 13 (running) created at:
  main.main()
      /Users/josetheodoro/workspace/learning/golang/concurrency/exercices/counter/main.go:15 +0x44

Goroutine 19 (finished) created at:
  main.main()
      /Users/josetheodoro/workspace/learning/golang/concurrency/exercices/counter/main.go:15 +0x44
==================
Todos os jobs finalizaram e counter= 948
Found 2 data race(s)
exit status 66
```

#### Fix Race Conditio`
A solução mais viável aqui é Mutex, do pacote sync. Mutex controla o acesso exclusivo do ciclo em execução a memória por meio de travas.

```go
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

```
Quando goroutine executa a função increment a primeira que coisa que acontece é o Lock, ou seja, travamos para que somente aquele ciclo possa ler e escrever em counter, após a operação de leitura e escrita, o Unlock é executado, liberando o acesso a outros ciclos.

Nota: Mutex faz calculos matemáticos para travar com eficiência.

Ao executarmos nosso código agora
```bash
❯ go run main.go
Todos os jobs finalizaram e counter= 1000
```

O resultado é 1000 certinho, sem nenhuma sobrescrita. Ao rodar com a flag -race

```bash
❯ go run -race main.go
Todos os jobs finalizaram e counter= 1000
```
Sem nenhum Warnging de race condition, resolvemos o problema de race condiiton.