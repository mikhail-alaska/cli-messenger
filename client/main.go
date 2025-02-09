package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {

	var token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzkyMjIzMjUsInVzZXJuYW1lIjoiYWxhc2thIn0.oip_l3FD2MZJJ0WcRchCn-qIsWvW9OW0lrULdlbJ5Xs"
	m := NewModel(token)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
