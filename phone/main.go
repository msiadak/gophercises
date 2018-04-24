package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type phoneNumber struct {
	ID          uint `gorm:"primary_key"`
	PhoneNumber string
}

func seedDB(db *gorm.DB) error {
	db.AutoMigrate(&phoneNumber{})

	seedNumbers := `
		1234567890
		123 456 7891
		(123) 456 7892
		(123) 456-7893
		123-456-7894
		123-456-7890
		1234567892
		(123)456-7892
	`

	var count uint
	db.Model(&phoneNumber{}).Count(&count)
	if count > 0 {
		return fmt.Errorf("Database already has seed data")
	}

	for _, number := range strings.Split(seedNumbers, "\n") {
		trimmed := strings.TrimSpace(number)
		if len(trimmed) > 0 {
			db.Create(&phoneNumber{PhoneNumber: strings.TrimSpace(number)})
		}
	}

	return nil
}

func normalizePhoneNumber(pn string) string {
	var buf bytes.Buffer

	for _, r := range pn {
		if unicode.IsNumber(r) {
			_, err := buf.WriteRune(r)
			if err != nil {
				panic(err)
			}
		}
	}

	return buf.String()
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %s\n", err)
	}
	defer db.Close()

	err = seedDB(db)
	if err != nil {
		fmt.Println(err)
	}

	phoneNumbers := make([]phoneNumber, 0)
	db.Find(&phoneNumbers)

	for _, rec := range phoneNumbers {
		rec.PhoneNumber = normalizePhoneNumber(rec.PhoneNumber)
		db.Save(&rec)

		var count uint
		db.Model(&phoneNumber{}).Where(&phoneNumber{PhoneNumber: rec.PhoneNumber}).Count(&count)

		if count > 1 {
			db.Delete(&rec)
		}
	}

	phoneNumbers = make([]phoneNumber, 0)
	db.Find(&phoneNumbers)
	fmt.Printf("Contents of database:\n")
	for _, rec := range phoneNumbers {
		fmt.Printf("%v\n", rec)
	}
}
