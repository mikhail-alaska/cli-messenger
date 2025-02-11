package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

func GetToken() string {
	user := cfg.Login
	pass := cfg.Password
	url := cfg.Address

	data := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, user, pass)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/login", url), strings.NewReader(data))
	if err != nil {
		return ""
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	gotBody := string(body)

	type TokenResponse struct {
		Status string   `json:"status"`
		Token  string `json:"Token"`
	}
	var response TokenResponse

	if err := json.Unmarshal([]byte(gotBody), &response); err != nil {
		log.Fatalf("GetToken: unable to parse find response: %v", err)
	}
	return response.Token
}
func GetUser() string {
	return cfg.Login
}

func getChats(jwtToken string) (string, error) {

	// Создаём новый запрос
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/chats", cfg.Address), nil)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/users", cfg.Address), nil)
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

func getuserfromitem(user item) string {
	return user.title
}

func sendMsg(selected list.Item, longMsg string) (bool, error) {
	rawuser := selected.(item)
	user := getuserfromitem(rawuser)
	rawMsgs := splitStringByMaxLen(longMsg, 0)
	jwtToken := GetToken()
	for _, rawMsg := range rawMsgs {
		msg := encodeMsg(rawMsg)

		data := fmt.Sprintf(`{"username": "%s", "messagefor1": "%s", "messagefor2": "%s"}`, user, msg, msg)
		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/message", cfg.Address), strings.NewReader(data))
		if err != nil {
			return false, err
		}
		req.Header.Set("Authorization", "Bearer "+jwtToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
	}

	return true, nil
}

func getMsg(user_raw item) (string, error) {
	var out string
	user := user_raw.title
	jwtToken := GetToken()
	data := fmt.Sprintf(`{"username": "%s"}`, user)
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/message", cfg.Address), strings.NewReader(data))
	if err != nil {
		return out, err
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}
	rawMsgs := string(body)

	return rawMsgs, nil
}
func splitStringByMaxLen(s string, maxLen int) []string {
	var result []string
	if maxLen == 0 {
		result := append(result, s)
		return result
	}
	runes := []rune(s)
	for len(runes) > 0 {
		if len(runes) > maxLen {
			result = append(result, string(runes[:maxLen]))
			runes = runes[maxLen:]
		} else {
			result = append(result, string(runes))
			break
		}
	}

	return result
}
