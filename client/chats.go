package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

func GetToken() string {

	var token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk0NzQyNDIsInVzZXJuYW1lIjoiYmVyb21uaWsifQ.cPmYiqjdgHfbuvIK_0cd2iktl_KwIv7fV8nFySn21Wo"
	return token
}
func GetUser() string {

	var user string = "alaska"
	return user
}

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

func getFind() (string, error) {
	// Создаём новый запрос
	req, err := http.NewRequest("GET", "http://localhost:8080/users", nil)
	if err != nil {
		return "", err
	}

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

func getuserfromitem(user item) string{
    return user.title
}

func sendMsg(selected list.Item, msg string) (bool, error) {
	rawuser := selected.(item)
    user := getuserfromitem(rawuser)

	jwtToken := GetToken()
	data := fmt.Sprintf(`{"username": "%s", "messagefor1": "%s", "messagefor2": "%s"}`, user, msg, msg)
	req, err := http.NewRequest("POST", "http://localhost:8080/message", strings.NewReader(data))
	if err != nil {
		return false, err
	}

	// Добавляем заголовок с JWT-токеном. Обычно его помещают в Authorization с префиксом Bearer.
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Используем клиент для отправки запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return true, nil
}

func getMsg(user_raw item) (string, error) {
    var out string
    user := user_raw.title
	jwtToken := GetToken()
	data := fmt.Sprintf(`{"username": "%s"}`, user)
	req, err := http.NewRequest("GET", "http://localhost:8080/message", strings.NewReader(data))
	if err != nil {
		return out, err
	}

	// Добавляем заголовок с JWT-токеном. Обычно его помещают в Authorization с префиксом Bearer.
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Используем клиент для отправки запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}
    rawMsgs := string(body)


	return rawMsgs, nil
}
