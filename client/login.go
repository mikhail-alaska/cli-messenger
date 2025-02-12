package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

func CheckRegister() string{

	user := cfg.Login
	url := cfg.Address

	data := fmt.Sprintf(`{"username": "%s"}`, user)
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/user", url), strings.NewReader(data))
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

	type UserResponse struct {
		Status string   `json:"status"`
		OpenKey  int `json:"OpenKey"`
	}
	var response UserResponse

	if err := json.Unmarshal([]byte(gotBody), &response); err != nil {
		log.Fatalf("GetToken: unable to parse find response: %v", err)
	}
    if response.Status == "Error" {
        ok := RegisterUser()
        if !ok {
            log.Fatalf("You can not be registered for some reason, contact the server owner")
        }
    }

    token := GetToken()
    return token
}

func RegisterUser() bool {
	user := cfg.Login
	pass := cfg.Password
    openkey := cfg.OpenKey
	url := cfg.Address

    data := fmt.Sprintf(`{"username": "%s", "password": "%s", "openkey": "%v"}`, user, pass, openkey)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/user", url), strings.NewReader(data))
	if err != nil {
		return false
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	gotBody := string(body)

	type RegisterResponse struct {
		Status string   `json:"status"`
	}
	var response RegisterResponse

	if err := json.Unmarshal([]byte(gotBody), &response); err != nil {
		log.Fatalf("Register: unable to parse find response: %v", err)
	}
    if response.Status == "Error"{
        return false
    }
	return true
}

func GetUser() string {
	return cfg.Login
}
