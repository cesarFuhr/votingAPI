# Voting microservice

Exposes an API to create discussions and register votes

# Observações
* Muito obrigado pela oportunidade, espero corresponder as expectativas
* Não fiz um fluxo de branch e merge pois foi um desenvolvimento sequencial e sem cooperação.
* O swagger com a documentação está na pasta /api
* Tem uma collection do postman também no /api
* Deixei os logs nas bordas para não gerar um log excessívo
* Qualquer dúvida estou a disposição

# Pontos propostos
* OK - Cadastrar uma nova pauta -> /agenda
* OK - Abrir uma sessão de votação -> /session
* OK - Receber votos associados a pautas -> /vote
* OK - Contabilizar os votos e dar o resultado da votação -> /result

# Pontos bonus
* OK - Integração com sistemas internos -> ver internal/app/adapters e internal/app/domain/vote
* OK - Mensageria e filas -> ver internal/app/adapters e internal/app/domain/session
* OK - Performance -> ver abaixo a descrição de como testar
* OK - Versionamento da API -> deixei um comentário em seguida, podemos conversar mais sobre o assunto

# Sobre versionamento da API
Em geral o aconselhável é que se utilize um versionamento logo no início do path da URL.

Um exemplo seria: /v1/agenda/...

É sempre bom ser retro compatível e permitir que o cliente tenha uma migração tranquila ou até use concomitantemente as APIs de versões diferentes.

Existem autores que não consideram correto o versionamento dentro dos microserviços (Susan Fowler). Segundo ela se entende o micro serviço quase como uma biblioteca errôneamente.
Cada mincro serviço tem seu ciclo de vida e se está obsoleto ou se atualiza o comportamento ou se reescreve do forma a ser mais aderente as novas regras de negócio.

# Como rodar

Para rodar a aplicação é necessário ter docker e docker-compose instalados.
Coloquei uma collection do Postman no /api para uso como exemplo.

## Build
Na raiz do projeto
```bash
docker-compose build
```

## Run
Na raiz do projeto
```bash
docker-compose up
```

## Stop and Clear
Na raiz do projeto
```bash
docker-compose down
```

## Acompanhar as métricas no dashboard MQTT
http://localhost:18083
User: admin
Pass: public

## Testar performance
Escolhi o endpoint de Result como o alvo do teste por ser o mais custoso sem contar o serviço de registro voto, que só é mais lento pois sai dos limites do serviço.
Deixei o teste com 2000 request/segundo, mas podem ficar a vontade para mudar tentar encontrar o teto em que o serviço retorna um 500.

### Buildar o serviço
Na raiz do projeto
```bash
docker-compose -f docker-compose.load.yml build
```
### Subir o serviço
Na raiz do projeto
```bash
docker-compose -f docker-compose.load.yml up -d app
```
### Com o postman (ou qulquer cliente http):
* POST /agenda
* POST /agenda/{agendaID}/session
* POST /agenda/{agendaID}/session/{sessionID}/vote (algumas vezes)

### Ajustar o docker-compose.load.yaml com o valores do agendaID e sessionID
* Linha 40 do arquivo

### Subir a aplicação de teste
```bash
docker-compose -f docker-compose.load.yml up test-runner
```

### Parar as aplicações e remover os containers
```bash
docker-compose -f docker-compose.load.yml down
```

### O resultado deve ser algo semelhante a este
```bash
test-runner_1  | Requests      [total, rate, throughput]         20000, 2000.11, 1762.35
test-runner_1  | Duration      [total, attack, wait]             11.349s, 9.999s, 1.349s
test-runner_1  | Latencies     [min, mean, 50, 90, 95, 99, max]  409.583µs, 91.236ms, 768.118µs, 271.479ms, 454.515ms, 1.33s, 3.555s
test-runner_1  | Bytes In      [total, mean]                     3020000, 151.00
test-runner_1  | Bytes Out     [total, mean]                     0, 0.00
test-runner_1  | Success       [ratio]                           100.00%
test-runner_1  | Status Codes  [code:count]                      200:20000
test-runner_1  | Error Set:
```