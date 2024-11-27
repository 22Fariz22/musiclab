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


## Features

### Получение данных библиотеки с фильтрацией по всем полям и пагинацией
- **Method**: `GetLibrary()`

### Получение текста песни с пагинацией по куплетам
- **Method**: `GetSongVerseByPage()`
- Добавлено кеширование в Redis, чтобы из-за каждого куплета не обращаться к PostgreSQL.

### Удаление песни
- **Method**: `DeleteSongByGroupAndTrack()`
- Удаление выполняется по названию группы и песни.  
  Однако, удаление по **ID** было бы более предпочтительным, безопасным и точным.

### Изменение данных песни
- **Method**: `UpdateTrackByID()`

### Добавление новой песни
- **Method**: `CreateTrack()`
- Для получения текста песни из стороннего API используется сервис:  
  [https://api.lyrics.ovh/v1/](https://api.lyrics.ovh/v1/)
- Для получения ссылки на YouTube и даты релиза используется **фиктивное значение**, так как бесплатного или не требующего ключа сервиса не найдено.
- Валидация текста песни **не выполняется**, так как песня может быть без слов.

### SWAGGER UI
- Доступен по адресу:  
  [https://localhost:8080/swagger/index.html](https://localhost:8080/swagger/index.html)
