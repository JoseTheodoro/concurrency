
##### A implementação Go: Modelo Híbrido (M:N)
O Go uniu o melhor dos dois mundos usando um modelo hibrido M:N. O runtime do Go cria um número M de Goroutines (que são threads leves a nível de usuário) e as distribuiu sobre um número N de kernel-level threads.
Dessa froma, o Go alcança paralelismo real usando multiplos núcleos, através das Kernel-threads, mas também a leveza e a velocidade criando milhões de goroutines (user-level) sobre elas.

###### LRQ e GRQ: Como o Go gerencial isso tudo?
Para orquestrar milhares ou milhões de goroutines rodando em cima das kernel-level threads, o Scheduler do Go utiliza das estruturas de fila de execução: LRQ e GRQ.

O Scheduler do Go trabalha com três entidades, o modelo GMP:
* **G Goroutine:** A tarefa que precisa ser executada.
* **M Machine:** A Kernel-level thread real do sistema operacinal
* **P Processor** Um contexto lógico do runtime do Go que gerencia as filas.

**LRQ - Local Run Queue**
Cada thread de Kernel recebe sua própia Local Run Queue
A LRQ mantém um subconjunto de goroutines que estão prontas para serem executadas por aquela thread específica.

Ter um fila local é vital para performance, se houvesse apenas uma única fila global para todas as threads, o Go precisaria usar um Mutex toda vez que uma thread quissesse pegar uma goroutine, criando um gargalo gigante. Com a LRQ, a thread pega a próxima goroutine de sua própria fila local rapidamente, sem precisar disputar recursos com outras threads.

**GRQ - Global Run Queue**
É uma fila central que armazena as goroutines que acabaram de ser criadas ou que ainda não foram atribuidas a nenhuma LRQ específica de uma thread do kernel. Periodicamente, as threads do kernel olham para essa fila buscando novos trabalhos.

**Work Stealing e Bloqueios de IO**
* **Work Stealing**: O runtime do Go usa essa mecanismo ou algoritmo distribuir o trabalho para que nenhuma thread fique ociosa. A thread que acabou o trabalho em sua LRQ, olha LRQ de outras threads buscando trabalho e pega pra si metade das goroutines. Também busca na GRQ buscando trabalho, isso garante que nenhum núcleo do processador fique parado enquanto houver trabalho pendente.
* **Bloqueios de IO:** Se uma goroutine em sua LRQ fizer uma operação demorada no sistema, operações de input/ouput, o SO pausaria a thread inteira e para evitar que todas as outras goroutines fiquem pressas aguardando, o Go intercpta essa chamada, desacopla toda LRQ para outra thread livre, criando uma nova se necessário, deixando a goroutine bloqueante em sua thread e LRQ original. as goroutines que ficariam bloqueadas, começam a execução imediatamente. Quando a goroutine terminar, ela é devolvida para LRQ existente ou para GRQ.

Goroutine que executa operação de IO em determinada Thread, o SO desprioriza essa thread e a deixa em estado de I/O waiting até que o que dispositivo retorne. (operações em disco e rede.) Nesse momento que o Go intercpta e faz o processo de transferencia para uma nova thread ou criando uma nova.

**Entendimento sobre operações de I/O** 
É toda comunicação do seu programa com componentes e dispositivos externos a CPU e memória principal. 
* Leitura e Escrita em disco
* Envio e Recebimento de mensagens pela Rede
* Espera de uma entrada do teclado do usuário.

Essas operações recebem muito destaque no estudo da concorrência porque costumam ser de ordens de grandeza mais lentas do que o processamento interno do CPU. Como a CPU precisaria ficar esperando os dispositivos fisicos, terminarem o trabalho, o sistema operacional e o runtime do Go aproveitam essas pausas causadas por I/O para colocar outras threads ou goroutines para trabalhar no processador, maximizando a eficiência do sistema.

**Pool de Kernel-Level Threads**
O runtime do Go determina quantas kernel-level threads ele vai usar para executar o código com base no número de processadores lógicos disponíveis no sistema.

A quantidade primária é definida e controloda por uma variável de ambiente chamada `GOMAXPROCS` . Se não configurar essa variável de forma manual, o runtime consultará o sistema operacional ao iniciar o programa para descobrir quantas CPUs a maquina tem e usurá esse valor como padrão.

**Scheduler do Go**
Tudo isso é orquestrado pelo Scheduler no runtime do Go. A criação, gerenciamento das threads de kernel, interceptação ao perceber que uma thread inteira será pausada pelo SO e provisionar uma outra thread, enquanto a thread aguarda na fila de espera para operações I/O a nível do SO, trabalhando juntamente com LRQ e GRQ. 