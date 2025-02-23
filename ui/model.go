package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"gorm.io/gorm"
)

type Model struct {
	List        list.Model
	Table       table.Model
	DB          *gorm.DB
	Message     string
	CurrentPage int
	FocusTable  bool
}
