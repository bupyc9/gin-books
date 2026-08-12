package database

import (
	"books/authors"
	"books/books"
	"books/config/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDb() *gorm.DB {
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  database.Dsn(),
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{},
	)
	if err != nil {
		panic("failed to connect database")
	}

	MigrateDb(db)

	return db
}

func MigrateDb(db *gorm.DB) {
	db.AutoMigrate(&authors.Author{})
	db.AutoMigrate(&books.Book{})
}
