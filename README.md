# Установка и запуск
1. установить утилиты: линтер и goose для миграций:
```bash
make install-golangci-lint
make install-goose
```

2. Добавить переменные окружения - ключ шифрования для токена:
```bash
export JWT_KEY="-my-256-bit-secret-"
```


# Инициализация базы данных
1. развернуть Docker контейнер с Postgres
```bash
docker-compose up -d
```
2. накатить миграции
```bash
make local-migration-up
```
3. запустить бэк
```bash
go mod init alfachatback
go mod tidy
go run cmd/chat/main.go
```
4. инициализировать базу через постман:
* В `Postman` импортировать коллекцию и окружение из папки `postman`.
* В `Postman` справа сверху выбрать импортированное окружение `alfachat`.

Из коллекции запустить: 
* `Admin/signup admin`
* `Admin/signin admin`
* `Admin/create system chat`

Таким образом создадим системного пользователя и системный чат.
Дальше можно добавлять пользователей, создавать чаты, добавлять в чаты пользователей,отправлять и получать сообщения.







# Инициализация базы данных
1. развернуть Docker контейнер с Postgres
```
docker-compose up
```
2. накатить миграции
```
make local-migration-up
```
3. запустить бэк
4. инициализировать базу через постман:
* `Admin/signup admin`
* `Admin/signin admin`
* `Admin/create system chat`

Таким образом создадим системного пользователя и системный чат.
Дальше можно добавлять пользователей, создавать чаты, добавлять в чаты пользователей,отправлять и получать сообщения.

5. 


# Установка и запуск
1. установить утилиты: линтер и goose для миграций:
```
make install-golangci-lint
make install-goose
```
2. В папку `bin` сгенерировать и положить файлы сертификата и ключа для протокола  https:
```
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```
3. Добавить переменные окружения - ключ шифрования для токена и пути нахождения файлов для http:
```
export TLS_CERT="/home/user/go/src/test-task/tasks/backend/GO/gremiha3/cert.pem"
export TLS_KEY="/home/user/go/src/test-task/tasks/backend/GO/gremiha3/key.pem"
export JWT_KEY="-my-256-bit-secret-"
```
