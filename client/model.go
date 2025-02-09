package main

import (
	"encoding/json"
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

const (
	listView uint = iota
	chatView
	newMessageView
    newChatView
)
type model struct {
	list list.Model
    state uint
textarea  textarea.Model
	textinput textinput.Model
}

func NewModel(token string) model{
	rawChats, err := getChats(token)
    items:=[]list.Item{}

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

    for _, i := range resp.Chats{
        items = append(items, item{title: i})
    }
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Список чатов"
    return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

