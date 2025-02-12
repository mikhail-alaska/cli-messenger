package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

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
