package tables

import "time"

type User struct {
	ID             uint      `gorm:"primarykey"`
	Name           string    `gorm:"size:255"`
	AccountCreated time.Time `gorm:"type:timestamp"`
}
