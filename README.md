# go-grpc-auth

Библиотека для аутентификации и авторизации пользователей с использованием gRPC и JWT (JSON Web Tokens). Это решение обеспечивает регистрацию, вход в систему и управление сессиями для вашего Go-приложения.

## Репозиторий с описанием ProtoBuff контракта: https://github.com/qu0ta/pet-proto

## Описание

`go-grpc-auth` предоставляет простой способ добавления функционала аутентификации и авторизации в ваше приложение на Go с использованием gRPC и JWT. Проект включает:

- Регистрацию пользователей.
- Вход в систему с получением JWT токена.
- Проверку и обновление токенов.
- Обработку ошибок аутентификации.

## Возможности

- Регистрация пользователей с уникальными email и паролем.
- Вход в систему с получением JWT токена для доступа к защищённым эндпоинтам.
- Защищённые эндпоинты с использованием middleware для валидации токенов.
- Простота интеграции с другими сервисами через gRPC.

## Установка

Для использования библиотеки в вашем Go проекте, выполните команду:

```bash
go get github.com/qu0ta/go-grpc-auth
```

## Использование

### 1. Регистрация пользователя

Пример запроса для регистрации нового пользователя:

```go
package main

import (
    ...
)

func main() {
	client := auth.NewAuthClient(conn) // conn - gRPC сессия

	req := &auth.RegisterRequest{
		Email:    "example@example.com",
		Password: "securepassword",
		AppId:    12345,
	}

	resp, err := client.Register(context.Background(), req)
	if err != nil {
		log.Fatalf("Error registering user: %v", err)
	}

	fmt.Printf("User registered with ID: %s\n", resp.GetUserId())
}
```

### 2. Вход в систему и получение токена
После регистрации можно выполнить вход и получить JWT токен для доступа к защищённым данным:

```go
package main

import (
    ...
)

func main() {
	client := auth.NewAuthClient(conn)

	req := &auth.LoginRequest{
		Email:    "example@example.com",
		Password: "securepassword",
	}

	resp, err := client.Login(context.Background(), req)
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}

	fmt.Printf("User logged in successfully. Token: %s\n", resp.GetToken())
}
```
### 3. Проверка токена
Вы можете использовать JWT токен для проверки доступа к защищённым эндпоинтам:

```go
package main

import (
	...
)

func main() {
	tokenString := "YOUR_JWT_TOKEN"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})

	if err != nil {
		log.Fatalf("Invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Fatal("Invalid token claims")
	}

	fmt.Println("Token claims:", claims)
}
```
## Настройка и конфигурация
Вы можете настроить ключи JWT, время жизни токенов и другие параметры через переменные окружения или конфигурационные файлы.

## Тестирование
Для запуска тестов используйте команду:

```bash
go test ./tests/
```

