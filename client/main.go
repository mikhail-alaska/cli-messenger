package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikhail-alaska/cli-messenger/client/config"
)

var cfg = config.MustLoad()

var globalToken = ""


func main() {
	globalToken = CheckRegister()
    if !cfg.Check(){
		fmt.Println("error while parsing your data")
		os.Exit(1)
    }
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
