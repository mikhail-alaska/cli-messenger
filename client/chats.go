package main

import (
	"io"
	"net/http"
)

func getChats(jwtToken string) (string, error) {
	// Создаём новый запрос
	req, err := http.NewRequest("GET", "http://localhost:8080/chats", nil)
	if err != nil {
		return "", err
	}

	// Добавляем заголовок с JWT-токеном. Обычно его помещают в Authorization с префиксом Bearer.
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Используем клиент для отправки запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
