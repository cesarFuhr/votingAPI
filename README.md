# Voting microservice

Exposes an API to create discussions and register votes


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

## Acompanhar as meétricas no dashboard MQTT
http://localhost:18083
User: admin
Pass: public

## Testar performance
Escolhi o ondepoint de Result como o alvo do teste por ser o mais custoso sem contar o serviço de registro voto, que só é mais lento pois sai dos limites do serviço.

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