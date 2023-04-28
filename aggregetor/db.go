package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func migrateDB(db *gorm.DB) {
	db.AutoMigrate(&Task{})
}

func initDBConnection() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=go_tasks port=5432 sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Can't Connect to postgresDB")
	}

	migrateDB(db)

	return db
}

func insertTasksToDB(tasks []Task) {
	db := initDBConnection()

	db.Create(tasks)

}
