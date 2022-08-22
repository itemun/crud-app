# CRUD-приложение для управления списком автомобилей

### Стэк
- go 1.19
- postgres 

### Запуск
```go build -o app cmd/main.go && ./app```

Для postgres можно использовать Docker

```docker run -d --name nina-db -e POSTGRES_PASSWORD=goLANGn1nja -v ${HOME}/pgdata/:/var/lib/postgresql/data -p 5433:5432 postgres```