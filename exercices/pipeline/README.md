## Pipeline em 3 estágios

Objetivo: Estruturar um encadeamento onde a saída de um processo assíncrono alimenta imediatamente a entrada de outro.

Vou codar  um fluxo contínuo de dados. 
1. Estágio 1: Codar um Gerador que envia os números de 1 a 20 para um canal. 
2. Estágio 2: Codar um Processador que lê do primeiro canal, multiplica o número por ele mesmo e envia para um segundo canal. 
3. Estágio 3: Codar um Consumidor que lê do segundo canal e imprime o valor no console. Os três estágios devem rodar concorrentemente.

> Todos os três estágios devem operar concorrentemente, processando a "esteira" de números em tempo real.