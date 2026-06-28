## Rate Limiter
Objetivo: Aprender a estrangular o processamento (throttling) utilizando time.Ticker, uma técnica vital para evitar o esgotamento de recursos e bloqueios ao consumir APIs externas com limites de requisição.

Exemplo do Worker Pool, independentemente de quantos trabalhadores existam na fila, o sistema como um todo só tem permissão para processar, no máximo, 2 tarefas por segundo.