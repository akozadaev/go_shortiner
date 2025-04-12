
# Проект [go_shortiner](https://github.com/akozadaev/go_shortiner) 
предназначен для демонстрации некоторых возможностей языка Go

## Показывает работу с БД и ORM

PostgreSQL + GORM, подумать над миграциями (авто)
## Показывает работу с фреймворками, зависимостями

Gin Web Framework, Fx

##  Журналирование, ротация логов, трассировка

Zap, Jaeger
![img.png](docs/img.png)

## Показывает работу с паттерном проектирования  Medaitor

Метод с общим названием Generate

## Для запуска

скачать репозиторий, создать пeстую базу Postgres ```link```, выполнить в консоли ```make build```, запустить собранное приложение.

### Роуты:
1. 
[POST] /v1/short

с телом объект или массив объектов вида:
```JSON
[
  {
    "url": "http://longlonglonglonglonglonglonglongurl.url"
  },
  {
    "url": "http://longlonglonglonglonglonglonglonglonglonglonglonglongurl.url"
  }
]
```
или
```JSON
{
  "url": "http://longlonglonglonglonglonglonglonglonglonglonglonglongurl.url"
}
```

2. 
[GET] /v1/short/:shortened

ответ вида:
```JSON
{
    "ID": 4,
    "CreatedAt": "2025-01-20T00:34:48.852775+03:00",
    "UpdatedAt": "2025-01-20T00:34:48.852775+03:00",
    "DeletedAt": null,
    "source": "http://urlEEEEEEEEee.url",
    "shortened": "https://short.ru/4339487037079594733"
}
```

Как может выглядеть вставка в таблицу для подготовки отчетов
```
INSERT INTO PREPARED_report (id, created_at, updated_at, timestamp, source, shortened, user_email, user_fullname)  
VALUES (1, '2025-03-15T12:00:00Z', now(), '2025-03-15T12:00:00Z', 'https://example.com', 'exmpl', 'user@example.com', 'John Doe');
```