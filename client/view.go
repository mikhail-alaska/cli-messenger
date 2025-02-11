package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle       = lipgloss.NewStyle().Margin(1, 2)
	myDescStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	otherDescStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	titleStyle     = lipgloss.NewStyle().Bold(true)
)

type itemDelegate struct {
	mydescStyle lipgloss.Style
	otherdescStyle lipgloss.Style
}

// Height возвращает высоту элемента (можно изменить, если нужен перенос строк).
func (d itemDelegate) Height() int { return 1 }

// Spacing между элементами.
func (d itemDelegate) Spacing() int { return 0 }

// Update обрабатывает сообщения (если не требуется — оставляем пустым).
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, itm list.Item) {
	i, ok := itm.(item)
	if !ok {
		return
	}

	// Определяем, выбран ли элемент.
	selected := index == m.Index()
	var renderedTitle string
	if selected {
		renderedTitle = titleStyle.Foreground(lipgloss.Color("229")).Render(i.Title())
	} else {
		renderedTitle = titleStyle.Render(i.Title())
	}

	// Выравнивание заголовка до фиксированной ширины.
	paddedTitle := lipgloss.NewStyle().Width(70).Render(renderedTitle)
    var renderedDesc string
	if i.Description() == GetUser() {
		renderedDesc = d.mydescStyle.Render(i.Description())
	} else {
		renderedDesc = d.otherdescStyle.Render(i.Description())
	}

	// Можно использовать разделитель или просто пробелы.
    s := fmt.Sprintf("%s:  %s", renderedDesc, paddedTitle)
	fmt.Fprint(w, s)
}
func (m model) View() string {
	switch m.state {
	case newMessageView:
		selected := m.currChat
		rawUser := selected.(item)
		user := getuserfromitem(rawUser)
		s := docStyle.Render(fmt.Sprintf("New message to %s", user)) + "\n\n"
		s += m.textinput.View() + "\n\n"
		return docStyle.Render(s)
	case listView:
		return docStyle.Render(m.list.View())
	case findView:
		return docStyle.Render(m.list.View())
	case chatView:
		delegate := itemDelegate{mydescStyle: myDescStyle, otherdescStyle: otherDescStyle}
		m.list.SetDelegate(delegate)
	}
	return docStyle.Render(m.list.View())
}
