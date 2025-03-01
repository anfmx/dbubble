package db

import (
	"fmt"
	"time"

	"github.com/Anfmx/dbubble/tables"
	"github.com/charmbracelet/bubbles/table"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var counter int

func ConnectDB(port, username, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=testdb port=%s sslmode=disable", username, password, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&tables.User{})
	return db, nil
}

func CreateTable(db *gorm.DB) error {
	if !db.Migrator().HasTable(&tables.User{}) {
		if err := db.Migrator().CreateTable(&tables.User{}); err != nil {
			return err
		}
	}
	return nil
}

func DropTable(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&tables.User{}); err != nil {
		return err
	}
	return nil
}

func CreateRow(db *gorm.DB) {
	counter++

	userExample := tables.User{Name: fmt.Sprintf("Test %d", counter), AccountCreated: time.Now()}
	db.Select(&tables.User{}).Create(&userExample)

	newUser := userExample
	newUser.ID = 0
	db.Create(&newUser)

}

func DeleteLastRow(db *gorm.DB) {
	db.Last(&tables.User{}).Delete(&tables.User{})
}

func GetColumns(db *gorm.DB) ([]table.Column, error) {
	var columns []table.Column

	rows, err := db.Raw(`
	SELECT column_name
	FROM information_schema.columns
	WHERE table_name = 'users' AND table_schema = 'public'`).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, err
		}
		columns = append(columns, table.Column{
			Title: columnName,
			Width: 15,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}
