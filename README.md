# Voting microservice

Exposes an API to create discussions and register votes

# Observações
* Não fiz um fluxo de brach e merge pois foi um desenvolvimento sequencial e sem cooperação.
* O swagger com a documentação está na pasta /api
* Deixei os logs nas bordas para não gerar um log excessívo
* Qualquer dúvida estou a disposição

# Sobre versionamento da API
Em geral o aconselhável é que se utilize um versionamento logo no início do path da URL.
É sempre bom ser retro compatível e permitir que o cliente tenha uma migração tranquila ou até use concomitantemente as APIs de versões diferentes.
Um exemplo seria: /v1/agenda/...
Existem autores que não consederam correto o versionamento detro dos microserviços (Susan Fowler). Segundo ela se entende o micro serviço quase como uma biblioteca errôneamente.
Cada mincro serviço tem seu ciclo de vida e se está obsoleto ou se atualiza o comportamento ou se reescreve do forma a ser mais aderente as novas regras de negócio.

# Como rodar

Para rodar a aplicação é necessário ter docker e docker-compose instalados.

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

## Acompanhar as meétricas no dashboard MQTT
http://localhost:18083
User: admin
Pass: public

## Testar performance
Escolhi o endpoint de Result como o alvo do teste por ser o mais custoso sem contar o serviço de registro voto, que só é mais lento pois sai dos limites do serviço.
Deixei o teste com 2000 request/segundo, mas podem ficar a vontade para mudar tentar encontrar o teto em que o serviço retorna um 500.

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