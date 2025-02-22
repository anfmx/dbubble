package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/username/tui-db-manager/bubbletea"
	"github.com/username/tui-db-manager/db"
)

// ANSI Styles
const (
	Reset     = "\033[0m"  // Styles reset
	Bold      = "\033[1m"  // Bold font
	WhiteText = "\033[37m" // White color
	RedBg     = "\033[41m" // Red background
	Hint      = "\033[34m" // Blue color
)

func main() {

	db, err := db.ConnectDB()
	if err != nil {
		fmt.Println(RedBg + Bold + WhiteText + "Error connecting to database\n" + Reset + Hint + "\nHint: try to enable postrges server in services" + Reset)
		fmt.Println(Bold + "\nPress enter to quit" + Reset)
		fmt.Scanln()
		os.Exit(1)
	}

	m := bubbletea.Model{
		CurrentPage: bubbletea.IntroPage,
		List:        list.New(bubbletea.CMDGreetingsPage, list.NewDefaultDelegate(), 0, 0),
		DB:          db,
	}
	m.List.Title = "Postresql TUI Manager"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
