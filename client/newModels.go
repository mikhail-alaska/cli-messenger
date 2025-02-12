package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func NewModel() model {

	token := GetToken()
	rawChats, err := getChats(token)
	items := []list.Item{}

	type ChatResponse struct {
		Status string   `json:"status"`
		Chats  []string `json:"Chats"`
	}
	var resp ChatResponse

	// Парсим JSON
	err = json.Unmarshal([]byte(rawChats), &resp)
	if err != nil {
		log.Fatalf("unable to parse chat response: %v", err)
	}

	ta := textarea.New()
	ti := textinput.New()

	for _, i := range resp.Chats {
		items = append(items, item{title: i})
	}
	if len(items) == 0 {
		items = append(items, item{title: "напиши кому нибудь"})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0), textarea: ta, textinput: ti}
	m.list.Title = "Список чатов"
	_ =m.list.AdditionalFullHelpKeys
    m.list.SetShowHelp(false) 
    m.list.SetStatusBarItemName("chat", "chats")
	return m
}

func ListModel(m *model) []tea.Cmd {
	token := GetToken()
	rawChats, err := getChats(token)
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	type ChatResponse struct {
		Status string   `json:"status"`
		Chats  []string `json:"Chats"`
	}
	var resp ChatResponse

	// Парсим JSON
	err = json.Unmarshal([]byte(rawChats), &resp)
	if err != nil {
		log.Fatalf("unable to parse chat response: %v", err)
	}

	newItems := []list.Item{}
	for _, chat := range resp.Chats {
		newItems = append(newItems, item{title: chat})
	}
	cmd = m.list.SetItems(newItems)
	cmds = append(cmds, cmd)
	m.list.Title = "Список чатов"
    m.list.SetStatusBarItemName("chat", "chats")
	return cmds
}

func ChatModel(m *model, selected list.Item) []tea.Cmd {
	user, ok := selected.(item)
	if !ok {
		log.Println("ChatModel: selected item имеет неверный тип")
	}

	rawChats, err := getMsg(user)
	if err != nil {
		log.Fatalf("ChatModel: unable to get messages: %v", err)
	}

	type MsgsResponse struct {
		Status string   `json:"status"`
		Msgs   []string `json:"Msgs"`
	}
	var resp MsgsResponse

	if err := json.Unmarshal([]byte(rawChats), &resp); err != nil {
		log.Fatalf("ChatModel: unable to parse chat response: %v", err)
	}

	var newItems []list.Item
	for _, msg := range resp.Msgs {
		parts := strings.SplitN(msg, ":", 2)
		if len(parts) == 2 {
            newItems = append(newItems, item{title: parts[1][2:], desc: parts[0]})
		} else{
			newItems = append(newItems, item{title: msg})
        }
	}

	m.list.SetItems(newItems)
	for range newItems {
		m.list.CursorDown()
	}
	m.list.Title = "Список сообщений"
    m.list.SetStatusBarItemName("message", "messages")
	return nil
}

func FindModel(m *model) []tea.Cmd {
	// Получаем данные по поиску
	rawChats, err := getFind()
	if err != nil {
		log.Fatalf("FindModel: unable to get find data: %v", err)
	}

	type FindResponse struct {
		Status string   `json:"status"`
		Users  []string `json:"Users"`
	}
	var resp FindResponse

	if err := json.Unmarshal([]byte(rawChats), &resp); err != nil {
		log.Fatalf("FindModel: unable to parse find response: %v", err)
	}

	var newItems []list.Item
	for _, user := range resp.Users {
		newItems = append(newItems, item{title: user})
	}

	m.list.SetItems(newItems)
	m.list.Title = "Список пользователей"
    m.list.SetStatusBarItemName("chat", "chats")
	return nil
}
