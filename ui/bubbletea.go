package ui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/username/tui-db-manager/db"
	"github.com/username/tui-db-manager/tables"
)

type item struct {
	title, desc string
}

const (
	IntroPage = iota
	MainPage
	TablePage
)

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var usersAmount int64
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "й":
			return m, tea.Quit
		case "enter":
			if m.CurrentPage == IntroPage {
				m.CurrentPage = MainPage
				m.List.Title = "DBubble"
				m.List.SetItems(CMDMainPage)
			} else if m.CurrentPage == MainPage {
				selectedItem, ok := m.List.SelectedItem().(item)
				if ok {
					switch selectedItem.Title() {
					case "Create table":
						if err := db.CreateTable(m.DB); err != nil {
							m.Message = errorStyle.Render("Error when creating table: " + err.Error())
						} else {
							m.Message = successStyle.Render("Table created successfully")
						}
					case "Drop table":
						if err := db.DropTable(m.DB); err != nil {
							m.Message = errorStyle.Render("Error when dropping table: " + err.Error())
						} else {
							m.Message = successStyle.Render("Table dropped successfully")
						}
					case "Insert test row":
						db.CreatRow(m.DB)
						m.Message = successStyle.Render("Row created")
						m.LoadTableData()
					case "Delete last row":
						if usersAmount > 0 {
							db.DeleteLastRow(m.DB)
							m.Message = successStyle.Render("Row created")
							m.DB.Model(tables.User{}).Count(&usersAmount)
						}
						if usersAmount == 0 {
							m.Message = notification.Render("There is no rows in table")
						}
						m.LoadTableData()

					case "Show data":
						m.LoadTableData()
						if len(m.Table.Rows()) == 0 {
							m.Message = notification.Render("No data to display")
						} else {
							m.Message = ""
						}
						m.CurrentPage = TablePage
						m.FocusTable = true
					}
				}
			}
		case "backspace":
			if m.CurrentPage == MainPage {
				m.CurrentPage = IntroPage
				m.List.Title = "Postresql TUI Manager"
				m.List.SetItems(CMDGreetingsPage)
			} else if m.CurrentPage == TablePage {
				m.CurrentPage = MainPage
				m.FocusTable = false
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
		if m.CurrentPage == TablePage {
			m.Table.SetWidth(msg.Width - h)
		}
	}

	if m.FocusTable && m.CurrentPage == TablePage {
		m.Table, cmd = m.Table.Update(msg)
	} else {
		m.List, cmd = m.List.Update(msg)
	}

	return m, cmd
}

func (m *Model) LoadTableData() {
	var users []tables.User
	res := m.DB.Find(&users)

	if res.Error != nil {
		m.Message = errorStyle.Render("Error fethcing data")
		return
	}
	if len(users) == 0 {
		m.Message = notification.Render("No data in table")
		m.Table.SetRows(nil)
		return
	}

	columnNames, err := db.GetColumns(m.DB)
	if err != nil {
		return
	}

	var rows []table.Row
	for _, user := range users {
		rows = append(rows, table.Row{
			strconv.Itoa(int(user.ID)),
			user.Name,
			user.AccountCreated.Format("2011-01-01 11:00")},
		)
	}

	m.Table = table.New(
		table.WithColumns(columnNames),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	m.Table.SetStyles(s)

	m.Table.SetRows(rows)
}

func (m Model) View() string {
	switch m.CurrentPage {
	case IntroPage:
		return docStyle.Render(m.List.View())
	case MainPage:
		return docStyle.Render(m.List.View() + "\n" + m.Message)
	case TablePage:
		return docStyle.Render(m.Table.View() + "\nPress [backspace] to return")
	default:
		return errorStyle.Render("Page does not exist")
	}
}
