package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	findView
)

type model struct {
	list      list.Model
	state     uint
	textarea  textarea.Model
	textinput textinput.Model

	viewport    viewport.Model
	messages    []string
	senderStyle lipgloss.Style
	err         error
	currChat    list.Item
}

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
	return cmds
}

func ChatModel(m *model, selected list.Item) []tea.Cmd {
	// Приводим выбранный элемент к нужному типу
	user, ok := selected.(item)
	if !ok {
		log.Println("ChatModel: selected item имеет неверный тип")
	}

	// Получаем данные сообщений для выбранного пользователя/чата
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

	// Создаём новый срез элементов для списка
	var newItems []list.Item
	for _, msg := range resp.Msgs {
		newItems = append(newItems, item{title: msg})
	}

	// Обновляем список сразу новым набором элементов
	m.list.SetItems(newItems)
	m.list.Title = "Список сообщений"
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
	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		key := msg.String()
		switch m.state {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "enter":
				m.state = chatView
                m.currChat = m.list.SelectedItem()
				cmds = append(cmds, ChatModel(&m, m.currChat)...)
				//cmds = append(cmds, tea.Println("user", m.list.SelectedItem()))
				//m.textarea.Focus()
				//m.textarea.CursorEnd()

			case "f", tea.KeyTab.String():
				m.state = findView
				cmds = append(cmds, FindModel(&m)...)
			}
		case findView:
			switch key {
			case "q":
				return m, tea.Quit
			case "enter", "esc":
				m.state = newMessageView
				m.currChat = m.list.SelectedItem()
				m.textinput.Focus()
				m.textinput.SetValue("")

			case tea.KeyTab.String():
				m.state = listView
				cmds = append(cmds, ListModel(&m)...)
			}

		case chatView:
			switch key {
			case "q":
				return m, tea.Quit

			case "u":
				cmds = append(cmds, ChatModel(&m, m.currChat)...)
			case tea.KeyTab.String():
				m.state = listView
				cmds = append(cmds, ListModel(&m)...)

			}
		case newMessageView:
			switch key {
			case "enter":
				msg := m.textinput.Value()
				if msg != "" {
					cmds = append(cmds, ChatModel(&m, m.list.SelectedItem())...)
					m.state = chatView
					sendMsg(m.currChat, msg)
				}
			case tea.KeyTab.String():
				m.state = listView
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	if m.state != newMessageView {
		m.list, cmd = m.list.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	item := m.list.SelectedItem()
	s := docStyle.Render(fmt.Sprintf("New message to %s", item)) + "\n\n"
	if m.state == newMessageView {
		s += m.textinput.View() + "\n\n"
		return docStyle.Render(s)
	}
	return docStyle.Render(m.list.View())
}
