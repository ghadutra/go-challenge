# Lastro - Go Challenge

## Descrição do Desafio

O objetivo deste desafio é implementar um sistema de chat utilizando Go e PostgreSQL, onde um assistente virtual responde às mensagens do usuário. O sistema deve processar mensagens de forma concorrente e garantir que as respostas do assistente sejam geradas após um atraso, permitindo que o usuário envie múltiplas mensagens antes de receber uma resposta.

### Descrição dos Arquivos

- **main.go:** Ponto de entrada do sistema, responsável por iniciar o servidor e as rotas.
- **service/routes.go:** Define as rotas HTTP para as operações de chat.
- **controllers/controller.go:** Implementa a lógica de negócios, incluindo o processamento de mensagens e a geração de respostas do assistente.
- **db/db.go:** Configura a conexão com o banco de dados PostgreSQL e define as operações de CRUD (Create, Read, Update, Delete) para chats e mensagens.

## Requisitos do Sistema

1. **Estrutura de Chat:**
    - Cada chat deve ser identificado por um ID único.
    - Cada chat pode ter múltiplas mensagens associadas a ele.
    - As mensagens devem indicar se foram enviadas pelo usuário ou pelo assistente.

2. **Comportamento de Resposta:**
    - O sistema deve esperar 10 segundos após receber uma mensagem antes de gerar uma resposta do assistente.
    - Se o usuário enviar várias mensagens dentro do intervalo de 10 segundos, o assistente deve responder a todas as mensagens recebidas em uma única resposta.
    - A resposta do assistente deve seguir o formato: `RESPOSTA DO ASSISTENTE PARA MENSAGEM: {X}`, onde `{X}` são as mensagens do usuário concatenadas.
    - **Atenção, não é necessário implementar nenhuma inteligência artificial ou chatbot. A resposta deve sempre seguir o padrão mencionado no ponto anterior**.
    - O sistema deve ser capaz de processar múltiplos chats simultaneamente, garantindo que as respostas sejam geradas corretamente.
3. **Banco de Dados:**
    - O sistema deve utilizar PostgreSQL para armazenar chats e mensagens.
    - As tabelas devem ser estruturadas para suportar múltiplas mensagens por chat, com indicações claras de autoria (usuário ou assistente).

4. **Testes:**
    - O sistema deve incluir testes para garantir que:
        - As mensagens são armazenadas corretamente no banco de dados.
        - O assistente responde corretamente após o intervalo de 10 segundos.
        - O comportamento de resposta em lote funciona conforme o esperado.

## Instruções para Implementação

1. **Implementação:**
    - Implemente a lógica de recebimento e processamento de mensagens no arquivo `controller.go`.
    - Configure o banco de dados PostgreSQL no arquivo `db.go` e garanta que as operações de CRUD funcionem corretamente.
    - Assegure que a lógica de espera de 10 segundos para a resposta do assistente esteja corretamente implementada.

2. **Testes:**
    - Escreva testes para verificar se as mensagens são corretamente armazenadas e recuperadas do banco de dados.
    - Teste se o assistente responde após o intervalo de 10 segundos e se responde corretamente a múltiplas mensagens.

3. **Entrega:**
    - Modifique o `README.md` com instruções sobre como rodar o sistema e os testes.
    - Submeta o repositório Git com a solução implementada.
    - Se existir SQL para criação do banco de dados, inclua-o em um arquivo `schema.sql`.

## Como Rodar

1. Clone o repositório:
    ```bash
       git clone <repo_url>
       cd go-challenge
    ```
2. Execute o sistema:
    ```bash
       go run main.go
    ```
   
## Testes

Para rodar os testes unitários e de integração:
```bash
go test ./...
```

## Benchmarks
Para executar os benchmarks:
```bash
go test -bench=.
```