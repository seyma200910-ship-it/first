# First — Go Backend Service

## 📌 Описание

Backend сервис на Go с использованием:

* PostgreSQL (pgx)
* Docker
* REST API

## ⚙️ Функционал

* Создание пользователя
* Получение пользователя по email
* Обновление пользователя
* Увольнение сотрудника (status = fired)

## 🧱 Архитектура

* repository layer
* service layer
* handler layer

## 🚀 Запуск

```bash
docker-compose up --build
```

## 📡 API

* `POST /users`
* `GET /users/{email}`
* `PATCH /users`
* `PATCH /users/{id}/terminate`

## 🛠 Технологии

* Go
* PostgreSQL
* Docker
