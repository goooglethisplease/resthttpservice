# REST-сервис подписок

Сервис агрегирует данные об онлайн-подписках пользователей, предоставляет CRUDL-операции и ручку для подсчета суммарной стоимости подписок за период. Используется PostgreSQL, миграции через Goose, документация Swagger.

## Возможности

- CRUDL для подписок (создание, чтение, обновление, удаление, список).
- Суммарная стоимость подписок за период с фильтрами.
- Логи в формате JSON.
- Конфигурация через `.env`.
- Swagger-документация.
- Запуск через `docker compose`.

## Модель подписки

- `service_name` — название сервиса
- `price` — стоимость месячной подписки в рублях (целое число)
- `user_id` — UUID пользователя
- `start_date` — дата начала в формате `MM-YYYY`
- `end_date` — опциональная дата окончания в формате `MM-YYYY`

### Создать подписку

`POST /api/subscriptions`

Пример тела запроса:

```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}
```

### Получить подписку

`GET /api/subscriptions/{id}`

### Список подписок

`GET /api/subscriptions?user_id=<uuid>&service_name=<name>`

### Обновить подписку

`PUT /api/subscriptions/{id}`

### Удалить подписку

`DELETE /api/subscriptions/{id}`

### Сумма подписок за период

`GET /api/subscriptions/summary?start_date=07-2025&end_date=09-2025&user_id=<uuid>&service_name=<name>`

Ответ:

```json
{
  "total": 1200
}
```
