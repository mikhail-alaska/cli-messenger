package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle        = lipgloss.NewStyle().Margin(1, 2)
	myDescStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	otherDescStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	titleStyle      = lipgloss.NewStyle().Bold(true)
	newMessageStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Margin(1, 2)
)

type itemDelegate struct {
	mydescStyle    lipgloss.Style
	otherdescStyle lipgloss.Style
}

func (d itemDelegate) Height() int { return 1 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, itm list.Item) {
	i, ok := itm.(item)
	if !ok {
		return
	}

	selected := index == m.Index()
	var renderedTitle string
	if selected {
		renderedTitle = titleStyle.Foreground(lipgloss.Color("229")).Render(i.Title())
	} else {
		renderedTitle = titleStyle.Render(i.Title())
	}

	paddedTitle := lipgloss.NewStyle().Width(m.Width() - 40).Render(renderedTitle)
	var renderedDesc string
	if i.Description() == GetUser() {
		renderedDesc = d.mydescStyle.Render(i.Description())
	} else {
		renderedDesc = d.otherdescStyle.Render(i.Description())
	}

	s := fmt.Sprintf("%s:  %s", renderedDesc, paddedTitle)
	fmt.Fprint(w, s)
}

func (m model) View() string {
	var mainContent string

	switch m.state {
	case newMessageView:
		selected := m.currChat
		m.textinput.Placeholder = "Start typing"
		rawUser := selected.(item)
		user := getuserfromitem(rawUser)
		title := fmt.Sprintf("New message to %s", user)
		titleStyled := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Render(title)
		inputView := m.textinput.View()
		content := fmt.Sprintf("%s\n\n%s", titleStyled, inputView)
		mainContent := newMessageStyle.Render(content)
		centeredContent := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, mainContent)
		return centeredContent
	case chatView:
		delegate := itemDelegate{mydescStyle: myDescStyle, otherdescStyle: otherDescStyle}
		m.list.SetDelegate(delegate)
		mainContent = docStyle.Render(m.list.View())
	default:
		mainContent = docStyle.Render(m.list.View())
	}

	helpMenuHeight := 5

	mainAreaHeight := m.height - helpMenuHeight

	mainContent = lipgloss.NewStyle().Height(mainAreaHeight).Render(mainContent)

	helpContent := renderHelpMenu(m.width)

	helpContent = lipgloss.Place(m.width, helpMenuHeight, lipgloss.Left, lipgloss.Bottom, helpContent)

	combined := lipgloss.JoinVertical(lipgloss.Left, mainContent, helpContent)
	return combined
}

func renderHelpMenu(screenWidth int) string {
	var helpItems []string

	helpItems = []string{
		"q, Ctrl+C: Quit (выход) | ↑/↓, k/j: Navigate (ориентироваться по списку)",
		"Tab/f: Change tab (сменить страницу) | Enter: Select (выбор) | / : Search (поиск) ",
        "g: Start of page | G: End of page | h/b: Previous page | l/d: Next page",
	}

	// Собираем список в одну строку с переносами.
	helpText := ""
	for _, item := range helpItems {
		helpText += item + "\n"
	}
	helpText = helpText[:len(helpText)-2]
	// Создаем стиль для меню помощи.
	// Убираем фиксированную ширину, поскольку мы будем задавать размеры через Place.
	helpMenuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0).
		Margin(0, 0)

	renderedHelp := helpMenuStyle.Render(helpText)

	// Используем lipgloss.Place, чтобы окно имело ровно 10 строк высотой и занимало всю ширину экрана.
	// Выравнивание по центру здесь применяется к содержимому внутри этого прямоугольника.
	finalHelp := lipgloss.Place(screenWidth, 2, lipgloss.Left, lipgloss.Center, renderedHelp)
	return finalHelp
}
