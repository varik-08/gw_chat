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
- `cp .env.example .env`
- `docker-compose up -d`

