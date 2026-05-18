# Trello Golang

Канбан-доска (аналог Trello) на Go.

## Возможности

- Регистрация и аутентификация (JWT)
- Управление досками, списками и карточками
- Комментарии к карточкам
- Прикрепление файлов к карточкам
- Публичные/приватные доски
- Загрузка пользовательских файлов

## Требования

- Go 1.25+
- PostgreSQL

## Установка и запуск

1. Клонируйте репозиторий:

```bash
git clone https://github.com/alexandervashurin/trello-golang.git
cd trello-golang
```

2. Скопируйте `.env` и настройте подключение к базе данных:

```env
DATABASE_URL=postgres://user:password@localhost:5432/trello_db
JWT_SECRET=your-secret-key
PORT=8080
```

3. Создайте базу данных PostgreSQL:

```bash
createdb trello_db
```

4. Запустите приложение:

```bash
go run main.go
```

Приложение будет доступно по адресу `http://localhost:8080`.

## API

### Аутентификация
- `POST /api/register` — регистрация
- `POST /api/login` — вход
- `GET /api/me` — текущий пользователь

### Доски
- `POST /api/boards` — создать доску
- `GET /api/boards` — все доски пользователя
- `GET /api/boards/public` — публичные доски
- `GET /api/board` — доска по ID

### Списки
- `POST /api/lists` — создать список
- `GET /api/lists` — списки доски

### Карточки
- `POST /api/cards` — создать карточку
- `GET /api/cards` — карточки списка
- `GET /api/card` — карточка по ID
- `PATCH /api/card` — переместить карточку
- `DELETE /api/card` — удалить карточку

### Комментарии
- `POST /api/comments` — создать комментарий
- `GET /api/comments` — комментарии карточки
- `DELETE /api/comment` — удалить комментарий

### Вложения
- `POST /api/attachments` — загрузить вложение
- `GET /api/attachments` — вложения карточки
- `DELETE /api/attachment` — удалить вложение
- `GET /api/files/{id}` — скачать файл
