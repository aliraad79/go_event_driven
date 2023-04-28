package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func migrateDB(db *gorm.DB) {
	db.AutoMigrate(&Task{})
}

func initDBConnection() *gorm.DB {

	if os.Getenv("DOCKER") == "false" {
		dsn := "host=localhost user=postgres password=postgres dbname=go_tasks port=5432 sslmode=disable TimeZone=Asia/Tehran"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Can't Connect to postgresDB")
		}
		migrateDB(db)

		return db
	} else {
		dsn := fmt.Sprintf("host=%s user=postgres password=postgres dbname=go_tasks port=5432 sslmode=disable TimeZone=Asia/Tehran", os.Getenv("POSTGRES_DOCKER_URL"))
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Can't Connect to postgresDB")
		}
		migrateDB(db)

		return db
	}

}

func insertTasksToDB(tasks []Task) {
	db := initDBConnection()

	db.Create(tasks)

}
