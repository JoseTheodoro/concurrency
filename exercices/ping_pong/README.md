## Ping-Pong
Objetivo: Praticar o envio e recebimento em canais bidirecionais sem causar interrupções ou deadlocks.

**Criar duas goroutines** 
* Goroutine 1 envia a palavra **"ping"** para um `channel`.
* Goroutine 2 recebe essa palavra, imprime no terminal e envia **"pong"** de volta.

**conferir a implementação em `main.go`**

### O que foi legal nesse exercício

#### Channels Unbuffered

Canais não bufferizados são síncronos. Significa que quando escrevemos uma mensagems no canal, é bloquenate até alguem ler. Essa é a regra que tem que ser respeitada e entendida.

```go
package main

func main (){
    ch := make(chan string)
    ch <- "mensagem"
}
```
Esse exemplo é bloqueante, a main trava em `ch <- "mensagem"` até um alguém ler do canal `ch`. Ao executar esse código, o erro Deadlock irá estourar.

```bash
fatal error: all goroutines are asleep - deadlock
```

Mesmo se adicionarmos a leitura do canal `ch`, o deadlock continua.
```go
package main

import "fmt"

func main (){
    ch := make(chan string)
    ch <- "mensagem"

    msg := <-ch
    fmt.Println("msg=", msg)
}
```
```bash
fatal error: all goroutines are asleep - deadlock
```

**Porque isso acontece?**
Canais foram projetados para comunicação e sincronização entre goroutines.
Todo programa go, é executado em pelo menos uma goroutine, a execução da `main` é em uma goroutine. O compilador do Go é inteligente suficiente para perceber que ao escrever `ch <- "mensagem"` e `ch` não é um canal buferizado ninguém vai conseguir ler a mensagem porque essa operação é bloqueante e um `deadlock` acontece.

**Como resolvemos?**
Acho que ficou óbivio, para `soltar o processo` precisamos apenas criar uma goroutine para escrever a  mensagem e não bloquear.
```go
package main

func main() {
    
    ch := make(chan string)
    go func(){
        ch <- "mensagem"
    }()
    
    msg:= <-ch
    fmt.Println("msg=", msg)
    time.Sleep(time.Second * 1) // apenas apra a main não terminar antes de ler o canal
}
```
Nesse caso é main aguarda a leitura da mensagem em `msg:= <-ch`, esse processo é bloquente até a mensagem chegar.
```bash
Output: msg= mensagem
```

**Adicioando outra Goroutine para ler**
```go
func main() {
	ch := make(chan string)

	go func() {
		ch <- "mensagem"
	}()

	go func() {
		msg := <-ch
		fmt.Println("msg=", msg)
		close(ch)
	}()

	time.Sleep(time.Second * 1)
}
```

**Adicionando boa prática de sinal de trabalho conclído**
Vamos adcioanar um cana para sinalizarmos quando o trabalho foi concluído.
`done := make(chan bool)`

```go
func main() {
	ch := make(chan string)
	done := make(chan bool)
	go func() {
		ch <- "mensagem"
	}()

	go func() {
		msg := <-ch
		fmt.Println("msg=", msg)
		close(ch)
		done <- true
	}()

	<-done
}
```
A Goroutine da main fica presa esperando alguma mensagem chegar no chanal done.
Ao terminar de ler a única mensagem na segunda goroutine, fechamos o canal para evitarmos memory leak e deadlocks e enviamos um sinal para o canal done.
```bash
Output:
msg= mensagem
main finalizada
```


#### Conslusão
No início foi dificil entender a ideia lógica por trás dos canais não buferizados, principalmente nessa orquestração de se alguém escreveu, outro tem que retirar. Fazer essa implementação foi bem legal para entender a teoria na prática. Fazer esse mergulo no modelo Actor ou CSP foi significativo para fixar o conhecimento. Exigiu de mim uma nova maneira de pensar, a leitura não sequencial dos acontecimentos e como orquestra-los.