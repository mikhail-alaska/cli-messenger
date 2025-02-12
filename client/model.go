package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

	width  int
	height int
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
			case tea.KeyCtrlC.String():
				return m, tea.Quit
			case "enter":
				m.state = chatView
				m.currChat = m.list.SelectedItem()
				cmds = append(cmds, ChatModel(&m, m.currChat)...)

			case "f", "l", "h", tea.KeyTab.String():
				m.state = findView
				cmds = append(cmds, FindModel(&m)...)
			}
		case findView:
			switch key {
			case tea.KeyCtrlC.String():
				return m, tea.Quit
			case "enter":
				m.state = newMessageView
				m.currChat = m.list.SelectedItem()
				m.textinput.Focus()
				m.textinput.SetValue("")

			case tea.KeyTab.String(), "f", "l", "h":
				m.state = listView
				cmds = append(cmds, ListModel(&m)...)
			}

		case chatView:
			switch key {
			case tea.KeyCtrlC.String():
				return m, tea.Quit

			case "u":
				cmds = append(cmds, ChatModel(&m, m.currChat)...)
			case "i":
				m.state = newMessageView
				m.textinput.Focus()
				m.textinput.SetValue("")
			case tea.KeyTab.String():
				m.state = listView
				cmds = append(cmds, ListModel(&m)...)

			}
		case newMessageView:
			switch key {

			case tea.KeyCtrlC.String():
				return m, tea.Quit
			case "enter":
				msg := m.textinput.Value()
				if msg != "" {
					sendMsg(m.currChat, msg)
					cmds = append(cmds, ChatModel(&m, m.currChat)...)
					m.state = chatView
				}
			case tea.KeyTab.String():
				m.state = listView
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	if m.state != newMessageView {
		m.list, cmd = m.list.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
