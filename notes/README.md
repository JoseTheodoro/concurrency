`goroutine` é o modelo básico de uma execução concorrente

Go nos dá muitas abstrações para coodenar execuções concorrentes em uma tarefa comum. Uma dessas abstrações é conhecidas como `channel`. Channel permite uma ou mais goroutine passar mensagens uma para outras, permitiando a troca de informação e sincronização de multiplas execuções de uma forma fácil e intuítiva.

Em 1978 foi descrito o modelo CSP _communicating sequencial processes_ para lidar com interações concorrentes. Muitas linguagens foram influenciadas por esse modelo, como Occam e Earlang. Go implementa muitas ideias da CSP, principalmente o uso de canais sincronizados `syncronized channels`.

Esse modelo de concorrencia com goroutines isoladas e canais reduz risco de race condition.

Dependendo do problema a programação concorrente clássica algumas vezes performam melhor que o estilo CSP. Go também prove algumas ferramentas para lidar com a programação classica.

Nesse livro vamos aprender como usar várias ferramentas para construir aplicações concorrentes. Isso inclui: Construções concorrentes como: mutex, varáveis de condição, canais, semaforos e assim por diante.

Go não vem com uma implementação padrão para semafaros, vamos implenentá-las do zero. Assim como outras implementações. A ideia central é conhecer essas ferramentas, tecnicas para construção de aplicações concorrentes mas também o entendimento de como elas podem trabalhar juntas.

#### Scaling Performance 
Escalabilidade de performance é a medição/indicador ou a forma de medir o quão bem a velocidade de um programa aumenta em proporção ao aumento do número de recursos disponíveis para o programa.
O livro passa a ideia que há um limite para a escalabilidade; Há um limite entre colocar mais rescursos e continuar diminuindo o tempo pela metade e apresenta `Amdahl's law`

##### Amdahl's law
Em 1967, Gene Amdahl, um cientista da computação, apresentou uma fórmula em uma conferência que mediu a aceleração em relação à relação paralelo-sequencial de um problema. Isso ficou conhecido como a lei de Amdahl.

**Definição** A lei de Amdahl afirma que a melhoria geral de desempenho obtida pela otimização de uma única parte de um sistema é limitada pela fração de tempo que a parte melhorada é realmente usada.

Que em outras palavras é: que partes não paralelas de uma de execução atuam como gargalos e limitam a vantagem de paralelizar a execução.

#### Go Runtime

##### Kernel-Level Threads
São threads que o sistema operacional gerencial através do kernel.

Sempre que executamos um programa, o sistema operacional gerencia os recursos e o tempo de CPU. O sistema operacional faz isso através de processose, dentro deles, as kernel-level threads
São as thread clássicas, como as primeiras versões do java, ou do C++ por exemplo. Elas são conhecidas, agendadas e gerenciadas internamente pelo sistema operacional.

**Vantagem** Se o processador tiver múltiplos núcleos de processamento, o sistema operacional pode colocar diferentes kernel-level threads rodando ao mesmo tempo em núcleos diferentes, entregando paralelismo real.

**Desvantagem** Elas são pesadas. Sempre que o sistema operacional precisa pausar uma thread para rodar outra, fazendo context switch, ele precisa intervir pesadamente. O SO precisa salvar o estado atual, limpar o cache do processador e carregar o estado da nova thread. Fazer isso milhares de vezes por segundo custa muita performance.

##### User-Level Threads
Para resolver o problema de lentidão das trocas de contexto, surgiram as User-level threads, frequetemente chamadas de _Green Threads_ em outras linguagens

São threads que rodam internamente dentro do espaço da meméria do programa. O sistema operacional não sabe que elas existem; Para o SO, o seu programa inteiro é apenas um única grande grande Kernel-level thread. Quem cria, pausa e gerencia essas threads não é mais o SO e sim o **runtime** da linguagem.

**Vantagem** Como o sistema operacional não está envolvido, a troca de contexto entre duas user-level threads é incrivelmente rápida. O runtime apenas troca os ponteiros de instrução na memória sem precisar acionar o núcleo do SO e limpar caches.

**Desvantagem** Como o SO enxerga apenas uma Kernel-level thread, se uma de suas user-level threads fizer uma operação demorada ou de bloqueio (operações de IO por exemplo), o SO pausa a kernel thread inteira. Isso significa que todas as outras user-level threads do seu programa também ficarão travadas, mesmo que estivessem prontas para rodar. Além disso, elas usarão apenas um núcleo do processador, desperdiçando o poder de máquinas com vários núcleos de processamento.

##### A implementação Go: Modelo Híbrido (M:N)
O Go uniu o melhor dos dois mundos usando um modelo hibrido M:N. O runtime do Go cria um número M de Goroutines (que são threads leves a nível de usuário) e as distribuiu sobre um número N de kernel-level threads.
Dessa froma, o Go alcança paralelismo real usando multiplos núcleos, através das Kernel-threads, mas também a leveza e a velocidade criando milhões de goroutines (user-level) sobre elas.

###### LRQ e GRQ: Como o Go gerencial isso tudo?
Para orquestrar milhares ou milhões de goroutines rodando em cima de um punhado de kernel-level threads, o Scheduler do Go utiliza das estruturas de fila de execução: LRQ e GRQ.

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