package bubbletea

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/username/tui-db-manager/db"
	"gorm.io/gorm"
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
	}
)

type item struct {
	title, desc string
}

const (
	IntroPage = iota
	MainPage
)

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Model struct {
	List        list.Model
	DB          *gorm.DB
	Message     string
	CurrentPage int
}

func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var usersAmount int64
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
					case "Delete last row":
						db.DeleteLastRow(m.DB)
						m.Message = successStyle.Render("Row deleted")
						m.DB.Model(db.User{}).Count(&usersAmount)
						if usersAmount == 0 {
							m.Message = notification.Render("There is no rows in table")
						}
					}
				}
			}
		case "backspace":
			if m.CurrentPage == MainPage {
				m.CurrentPage = IntroPage
				m.List.Title = "Postresql TUI Manager"
				m.List.SetItems(CMDGreetingsPage)
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.CurrentPage {
	case IntroPage:
		return docStyle.Render(m.List.View())
	case MainPage:
		return docStyle.Render(m.List.View() + "\n" + m.Message)
	default:
		return errorStyle.Render("What are you doing there?")
	}
}
