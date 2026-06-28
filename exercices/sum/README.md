## Soma Dividida

#### Objetivo
Entender como fracionar carga de trabalho e coletar resultados de execuções paralelas.

#### Problema
1. Criar um slice contendo 100 números inteiros aleatórios. 
2. Divida esse slice exatamente na metade. 
3. Inicie duas goroutines, enviando metade do slice para cada uma. Cada goroutine deve somar a sua parte e enviar o valor total para um canal numérico. 
4. A função main deve ler os dois valores desse canal, somá-los e exibir o resultado final.

#### O que foi legal nesse exercício

A implementação pode ser vista [main.go](main.go)

Eu não quis usar WaitGroup justamente para exercitar a pensar em como sincronizar as duas goroutines
WaitGroup é bem indicado para cargas de trabalho dinamicas, em meu cenário a carga de trabalho é estática, fixa, preferi não usar.

**Desafio de sincronização**
Acho que teria várias abordagens aqui de sincronização. o tempo de finalização de trabalho das duas goroutines é indeterminado. Como saber que elas terminaram seu trabalho e o processo da main possa computar o resultado final?

Em nosso cenário, sabemos exatamente o número de carga de trabalho e sabemos também que ela é imutável, fixa. Pela simplificação escolhi ler o canal duas vezes.
```go
var total int
for i:=1; i<=2; i++ {
    total = <-ch
}
```
Poderia sincronizar em uma linha também, sabendo que teremos sempre duas mensagens
```go
var total int
total += <-ch + <-ch
```
 **Especificamente para esse desafio** usar for lendo o canal não é indicado.
 ```go
var total int
for v := range ch {
    total += v
}
 ```
 For range espera mensagensinfinitamente quando usado para iterar canais(a menos que o canal seja fechado) e um erro de deadlock estouraria. Ficaria preso eternamente esperando a terceira mensagem.

 Nota: `select{}` não é indicado e nem é o propósito dele para iterar canais e somar resultados. É super indicado para multiplexação de canais, ou seja, escutar multiplos canais simultaneamente e agir no primeiro evento que chegar.