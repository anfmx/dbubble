package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss Styles
var (
	docStyle     = lipgloss.NewStyle().Margin(1, 2)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	notification = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

// Styled pages and titles
var (
	CMDGreetingsPage = []list.Item{
		item{title: "Greetings", desc: "This app creates and deletes tables in Postgres database"},
	}

	CMDMainPage = []list.Item{
		item{title: "Create table", desc: "Creating default table"},
		item{title: "Drop table", desc: "Deletes table"},
		item{title: "Insert test row", desc: "Inserts basic row"},
		item{title: "Delete last row", desc: "Deletes last created row"},
		item{title: "Show data", desc: "Display table contents"},
	}
)

// ANSI Styles
const (
	Reset     = "\033[0m"  // Styles reset
	Bold      = "\033[1m"  // Bold font
	WhiteText = "\033[37m" // White color
	RedBg     = "\033[41m" // Red background
	Hint      = "\033[34m" // Blue color
)
