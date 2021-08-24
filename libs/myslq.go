package libs

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (c *DbConfig) InitMysqlDB() *gorm.DB {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", c.User, c.Password, c.Host, c.Port, c.Database, c.Charset)
	db, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}

	sqlDB, err := db.DB()

	if err != nil {
		fmt.Println("Hubo un error", err)
		os.Exit(1)
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)

	return db
}
