# GW CHAT

[![Go Report Card](https://goreportcard.com/badge/github.com/varik-08/gw_chat)](https://goreportcard.com/report/github.com/varik-08/gw_chat)
[![CI Status](https://github.com/varik-08/gw_chat/actions/workflows/ci.yml/badge.svg)](https://github.com/varik-08/gw_chat/actions)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.14-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

<img src="web/chat/public/logo.jpg" alt="logo" width="200">

## Описание функционала

Чат-приложение для отправки сообщений в реальном времени.
С функционалом:
- Регистрации/Входа
- Личный кабинет с возможностью смены пароля и общей информацией о пользователе
- Личная переписка с отдельным пользователем
- Групповая переписка

## Технологии

- Golang
- WebSockets
- PostgreSQL
- Docker
- Docker Compose

## Развертывание
Развертывание сервиса должно осуществляться с использованием docker compose в директории с проектом.
- `cp .env.example .env` - копирование файла .env.example в .env
- `docker-compose up -d` - запуск сервиса
- `docker cp db/schema.sql <container_name>:/path/in/container/schema.sql` - копирование скрипта создания базы данных в контейнер
- `docker exec -i <container_name> psql -U <username> -d <database_name> -f /path/in/container/schema.sql` - выполнение скрипта создания базы данных в контейнере

## Экраны приложения
<div style="display: flex; flex-wrap: wrap;">
  <div style="margin: 10px;">
    <img src="docs/images/login.png" alt="Вход" width="300" />
  </div>
  <div style="margin: 10px;">
    <img src="docs/images/chat_list.png" alt="Список чатов" width="300" />
  </div>
<div style="margin: 10px;">
    <img src="docs/images/create_chat.png" alt="Создание чата" width="300" />
  </div>
<div style="margin: 10px;">
    <img src="docs/images/text.png" alt="Переписка" width="300" />
  </div>
  <div style="margin: 10px;">
    <img src="docs/images/profile.png" alt="Страница профиля" width="300" />
  </div>
</div>