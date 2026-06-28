## Channels

### Introdução
Como o design da linguagem Go foi inspirada no [modelo CSP](./csp_model.md) para lidar com concorrência, chegamos em channels, que é justamente para troca de mensagens entre processos concorrentes.

### Canais Não Bufferizados
Uma implementação mais hands-on da uma ideia prática sobre canais não bufferizados. [Ping-Pong](../exercices/ping_pong/README.md)

**Unbeffered Channels** São canais sem especificar um tamanho de `buffer`. Desse modo são **`estrimamente síncronos`**. Qualquer mensagem escrita nesse canal é bloqueante, a execução do programa é presa nesse momento e só é liberada quando outra goroutine consome, retira essa mensagem do canal do outro lado. O inverso também se aplica, quem tenta ler desse canal vazio fica travado até que alguém escreva.

> Quem é mais antigo, a melhor analogia é o telefone fixo: Pessoa A liga para Pessoa B, Pessoa A está pressa enquanto pessoa B não atende, quando Pessoa B atende, acontece o que chamamos de `curto` no `lic` e as duas se falam, trocam mensagens.

Em canais não bufferizados é o mesmo comportamente, precisa existir esse `curto`: goroutine A e B precisam se encontrar simultaneamente no canal para que a troca de mensagem aconteça.


### Canais Bufferizados

**Buffereds Channels**
São canais que são criados com um tamanho de `buffer` que representam a sua capacidade. Diferente dos canais não bufferizados, leitura e escrita não são bloqueantes, a menos que a sua capacidade, `buffer` esteja cheia, um channel bufferizado cheio tem o mesmo comportamento de canais não bufferizados.O comportamento de leitura é padrão de `queue`, FIFO (First-In, First-Out) mantida na memória do runtime do Go.

Nota: Internamente canais usam estruturas de semaforos de controle. Um semáfaro para controlar a capacidade e outro para controlar as mensagens recebidas. São mutexes para controlar o número de execuções concorrentes, Enquanto Mutex nos permite controlar apenas uma execução naquele momento.