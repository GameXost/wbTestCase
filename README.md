# Демонстрационный сервис с Kafka, PostgreSQL, кэшем

---
Сервис получает данные заказов из Kafka, сохраняет их в Postgres, а также кэширует в памяти и предоставляет HTTP API с простым интерфейсом для получения информации о заказе по его ID.

---

### Использовал:
- Kafka (franz-go)
- PostgreSQL (pgx)
- LRU cache (map + mutex)
- HTTP (chi)

### Структура:
/cache:
    cache.go
```plaintext
/cmd:
    main.go
    /producer
        main.go  - продюсер кафки

/config:
    config.go

./internal:
    /errHandle
        Errors.go
        
    /repository
        repo_bruh.go

    /server
        server_bruh.go

    /service
        service_bruh.go

/kafka:
    consumer.go

/models:
    delivery.go
    items.go
    orders.go
    payment.go

/web:
    index.html
```
---

## Запуск

- docker-compose up --build

### для тестирования кафки, можно запустить продюсер
    из корневой папки проекта выполнить
    go run ./cmd/producer/main.go 
---

##### По базе HTTP API доступен на 
`http://localhost:8080`
##### Еще прикручен UI для кафки
`http://localhost:8090`

##### Для получения информации по заказу доступен:
`GET /order/{order_uid}`
###### Также присутствует .env с переменными окружения, которые подтягиваются в main.go


#### P.S. работу кэша проверил
    2026/02/10 15:52:47 "GET http://localhost:8080/order/order_1770738456_14 HTTP/1.1" from 172.20.0.1:57210 - 200 988B in 2.850163ms
    2026/02/10 15:52:49 cache hit
    2026/02/10 15:52:49 "GET http://localhost:8080/order/order_1770738456_14 HTTP/1.1" from 172.20.0.1:57210 - 200 988B in 54.145µs
    2026/02/10 15:52:49 "GET http://localhost:8080/order/order_1770738456_14 HTTP/1.1" from 172.20.0.1:57210 - 200 988B in 82.2µs
