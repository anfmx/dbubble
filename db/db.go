package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var counter int

type User struct {
	ID             uint `gorm:"primarykey"`
	Name           string
	AccountCreated time.Time
}

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("env file loading error")
	}

	port := os.Getenv("PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=testdb port=%s sslmode=disable", dbUser, dbPass, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{})
	return db, nil
}

func CreateTable(db *gorm.DB) error {
	if !db.Migrator().HasTable(&User{}) {
		if err := db.Migrator().CreateTable(&User{}); err != nil {
			return err
		}
	}
	return nil
}

func DropTable(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&User{}); err != nil {
		return err
	}
	return nil
}

func CreatRow(db *gorm.DB) {
	counter++

	userExample := User{Name: fmt.Sprintf("Test %d", counter)}
	db.Select(&User{}).Create(&userExample)

	newUser := userExample
	newUser.ID = 0
	db.Create(&newUser)

}

func DeleteLastRow(db *gorm.DB) {
	db.Last(&User{}).Delete(&User{})
}
