# Birthday notifier service

Сервис позволяет зарегистрироваться, получить список пользователей, подписаться (и отписаться) на
интересующих пользователей и получать на почту уведомления об их днях рождения.

## Конфигурация

### Переменные окружения

- GOOSE_DRIVER - драйвер базы данных для `goose` - утилиты, используемой для миграций
- GOOSE_DBSTRING - строка подключения к базе данных для `goose`
- ENV - окружение, на котором запускается сервис
- APP_HOST - хост сервиса
- APP_PORT - порт сервиса
- JWT_SIGNING_KEY - ключ для подписи JWT
- JWT_TOKEN_TTL - время жизни выдаваемых JWT
- MAILER_EMAIL - email, используемый для рассылки уведомлений
- MAILER_PASSWORD - пароль для доступа к email'у выше
- MAILER_SMTP_HOST - хост SMTP-сервера используемого email'а
- MAILER_SMTP_PORT - порт SMTP-сервера
- MAILER_WAIT_BEFORE_RETRY - время ожидания перед следующей попыткой отправки письма, если была ошибка
- MAILER_MAX_RETRIES_COUNT - максимальное количество попыток отправки письма
- MAILER_INCREMENTAL_WAIT - увеличивать ли время ожидания в зависимости от номера попытки.
Например, если передать значение `true`, и если `MAILER_WAIT_BEFORE_RETRY` - 10 секунд, а `MAILER_MAX_RETRIES_COUNT` - 5,
то в случае ошибки на первой повторной попытке ожидание составит 10 секунд, на второй 20, и так далее
- DB_USER - имя пользователя в базе данных
- DB_PASSWORD - пароль пользователя в базе данных
- DB_HOST - хост базы данных
- DB_PORT - порт базы данных
- DB_NAME - название базы данных

Файл `.env.example` содержит все вышеперечисленные переменные с предзаполненными значениями.

## Запуск сервиса

Указать в файле `.env` необходимые переменные окружения.

```shell
docker network create bday-notifier
```

```shell
docker compose up
```

После запуска сервиса, по адресу `<APP_HOST>:<APP_PORT>/docs/` будет доступна swagger-документация.