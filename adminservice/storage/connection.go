package storage

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Psql(dbrole, dbpass, dbname, dbhost, dbport string, logger logger.Interface) *gorm.DB {
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN: fmt.Sprintf("user=%s password=%s dbname=%s host=%v port=%s sslmode=disable TimeZone=Asia/Shanghai",
				dbrole, dbpass, dbname, dbhost, dbport),
		}), &gorm.Config{Logger: logger})

	if err != nil {
		log.Fatalln("Database connection error")
	}

	return db
}
