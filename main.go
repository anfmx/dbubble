package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Anfmx/dbubble/db"
	"github.com/Anfmx/dbubble/ui"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var port string
	var username string
	var password string
	if _, err := os.Stat(".config"); os.IsNotExist(err) {
		fmt.Println("Database connection")

		fmt.Println("Port:")
		fmt.Scanln(&port)

		fmt.Println("Username:")
		fmt.Scanln(&username)

		fmt.Println("Password:")
		fmt.Scanln(&password)

		configData := fmt.Sprintf("Port: %s\nUsername: %s\nPassword: %s", port, username, password)

		file, err := os.Create(".config")
		if err != nil {
			os.Exit(1)
		}

		defer file.Close()

		if _, err := file.WriteString(configData); err != nil {
			fmt.Println("File writing error:", err)
			os.Exit(1)
		}
	}

	config := make(map[string]string)

	file, err := os.Open(".config")
	if err != nil {
		fmt.Println("Open file error:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			config[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("File reading error:", err)
		os.Exit(1)
	}

	database, err := db.ConnectDB(config["Port"], config["Username"], config["Password"])
	if err != nil {
		fmt.Println(ui.RedBg + ui.Bold + ui.WhiteText + "Error connecting to database\n" + ui.Reset + ui.Hint + "\nHint: try to enable postrges server in services" + ui.Reset)
		fmt.Println(ui.Bold + "\nPress enter to quit" + ui.Reset)
		fmt.Scanln()
		os.Exit(1)
	}

	columns, err := db.GetColumns(database)
	if err != nil {
		fmt.Println("Error fethcing column names:", err)
		os.Exit(1)
	}

	listModel := CreateListModel()
	tableModel := CreateTableModel(columns)

	m := ui.Model{
		CurrentPage: ui.IntroPage,
		List:        listModel,
		Table:       tableModel,
		DB:          database,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func CreateListModel() list.Model {
	listModel := list.New(ui.CMDGreetingsPage, list.NewDefaultDelegate(), 0, 0)
	listModel.Title = "Postgres TUI Manager"
	return listModel
}

func CreateTableModel(columns []table.Column) table.Model {
	return table.New(
		table.WithColumns(columns),
		table.WithHeight(7),
		table.WithFocused(true),
	)
}
