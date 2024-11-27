
# Music Library Example

  

## Full List of Technologies Used

- **[Echo](https://echo.labstack.com/)** - Web framework

- **[sqlx](https://github.com/jmoiron/sqlx)** - Extensions to `database/sql`

- **[pgx](https://github.com/jackc/pgx)** - PostgreSQL driver and toolkit for Go

- **[go-redis](https://github.com/redis/go-redis)** - Type-safe Redis client for Golang

- **[zap](https://github.com/uber-go/zap)** - Logger

- **[validator](https://github.com/go-playground/validator)** - Go Struct and Field validation

- **[swag](https://github.com/swaggo/swag)** - Swagger

- **[testify](https://github.com/stretchr/testify)** - Testing toolkit

- **[CompileDaemon](https://github.com/githubnemo/CompileDaemon)** - Compile daemon for Go

- **[Docker](https://www.docker.com/)** - Docker containerization

  

---

  

## Running the Application Locally

To start the application locally, run:

```bash

make up

```

  

## Features

  

### Получение данных библиотеки

- **Метод**: `GetLibrary()`

- Фильтрация по всем полям и поддержка пагинации.

  

### Получение текста песни

- **Метод**: `GetSongVerseByPage()`

- Пагинация по куплетам.

- Кеширование в Redis, чтобы минимизировать запросы в PostgreSQL.

  

### Удаление песни

- **Метод**: `DeleteSongByGroupAndTrack()`

- Удаление выполняется по названию группы и песни.

⚠️ Рекомендуется поменять на удаление по **ID** для большей точности и безопасности.

  

### Изменение данных песни

- **Метод**: `UpdateTrackByID()`

- Изменение данных трека по его ID.

  

### Добавление новой песни

- **Метод**: `CreateTrack()`

- Используется сторонний сервис для получения текста песни:

[https://api.lyrics.ovh/v1/](https://api.lyrics.ovh/v1/)

- Для заполнения полей ссылки и даты релиза используются фиктивные данные.

- Текст песни **не валидируется**, так как песня может быть без слов.

  

### Swagger UI

- Доступен по адресу:

[https://localhost:8080/swagger/index.html](https://localhost:8080/swagger/index.html)