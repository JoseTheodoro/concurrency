## Fundamentos da Concorrencia
Para um entendimento mais profundo, precisamos entender alguns pilares da concorrência.

#### Modelos de Comunicação 
A forma ou a estratégia de como vamos lidar com o fluxo de dados. Como as execuções irão trocar informações e sincroniza-las. Temos dois principais modelos de comunicação:
1. Clássico
    - É o modelo clássico onde threads compartilham o mesmo espaço de memória. Para eviar que uma thread altere o estado da outra no mesmo instante, conhecido como race condition é obrigatório o uso de travas como: Mutexes, Semáfaros e Variáveis de Consição. As primeiras linguagens de programação adotaram esse modelo de comunicação. Quem é das antigas vai lembrar que no processo do programa era comum rodar um `Fork()` para abrir uma thread compartilhando o mesmo espaço de memória que seu programa.
2. CSP
    - Modelo proposto em 1978 pelo cientista da computação Hoare. A premissa é a imutabilidade e o isolamento. As execuções não compartilham memória. Cada processo mantém o seu estado isolado, e quando quer compartilhar dados, envia `mensagens` através de `canais`. Linguagens como Erlang, Clojure, Go foram unfluenciadas por esse modelo.

#### Modelos de Execução e Agendamento
Definem e orquestram a execução dos processos nas unidades físicas ou lógica dos processadores. Orquestram o tráfego. Pelo que pude pesquisar temos 4 modelos de execução e agendamento:
1. Processos
    - Um processo é um programa em execução no sistema operacional. Ele possui um espaço de memória totalmente isolado, mas o custo de criação é alto e o consumo de memória também é alto. A comunicação entre processos (IPC) é complexa, exigindo manipular diretamente ferramentas do sistema operacional.
2. Kernel-Level Threads
    - Para reduzir o peso dos processos, surgiram as _threads_. Uma thread é uma linha de execução que vive dentro de um processo, compartilhando o mesmo espaço de memória, mas mantendo sua própria pilha de execução, registrdores e contador de instrução. São criadas e gerenciadas pelo sistema operacional, permitindo paralelismo real, cada thread pode rodar em núcleo do processador. O problema que ao fazer context-switch, a troca de contexto entre elas é custosa porque o sistema operacional precisa intervir, salvar estados na memória e limpar o cache do processador.
3. User-Level Threads
    - As famosas _Green Thread_. São threads gerenciadas internamente na memória a nível do usuário, runtime do programa em execução, sem que o sistema operacional saibam da sua existência. Para o SO existe apenas uma grande Kernel-Level Thread em execução. A troca de contexto é extremamente mais rápida, mas o grande problema  é: se uma green thread faz uma chamada de I/O bloqueante, o SO pausa a Kernel-Level Thread inteira, travando todas as green threads que estavam nela, não bastando isso, green threads também não conseguem paralelismo real justamente pelo seu design.
4. Corrotinas e o Modelo Híbrido M:N
 - Corrotinas são sub-função concorrentes que não são preemptivas, diferentemente das threads clássicas do sistema operacional, as corrotinas não podem ser suspensas pelo SO enquanto processam alguma tarefa. Elas possuem multiplos pontos ao longo de sua estrutura que permitem que a execução seja pausada e posteriormente retomadas. O go nos fornece uma abstração para executação de uma corrotina o comando ou palavra-chave `go`
 - O modelo hibrido M:N que o Go implementa: Go não usa threads do SO nem green threads puras, Usa justamente essas abstrações de corrotinas. Para executa-las o Go usa o modelo Híbrido M:N, que mapeia M goroutines (user-level threads) sobre N (kernal-level threads), garantindo paralismo real mesmo com execuções I/O. O runtine do Go faz a orquestrção em tres engrenagens: M: Thread nível do sistema operaional P: Processador, lógico ou físico  G: Corrotina. Como o `Scheduler` do Go faz essa orguestração, resolvi separar em [Go runtime](./go_runtime.md)