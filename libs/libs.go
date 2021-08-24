package libs

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/gorm"
)

var DB *gorm.DB

type DbConfig struct {
	Host         string
	Port         string
	Database     string
	User         string
	Password     string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
	TimeZone     string
}

func Configure(Env map[string]string) DbConfig {

	if Env["HOST"] == "" {
		fmt.Println("falta HOST en el archivo ")
		os.Exit(1)
	}

	if Env["DATABASE"] == "" {
		fmt.Println("falta DATABASE en el archivo " + "")
		os.Exit(1)
	}

	if Env["USER"] == "" {
		fmt.Println("falta USER en el archivo " + "")
		os.Exit(1)
	}

	if Env["PASSWORD"] == "" {
		fmt.Println("falta PASSWORD en el archivo " + "")
		os.Exit(1)
	}

	MAXIDLECONNS, _ := strconv.Atoi(Env["MAXIDLECONNS"])
	MAXOPENCONNS, _ := strconv.Atoi(Env["MAXOPENCONNS"])

	response := DbConfig{
		Host:         Env["HOST"],
		Port:         Env["PORT"],
		Database:     Env["DATABASE"],
		User:         Env["USER"],
		Password:     Env["PASSWORD"],
		Charset:      Env["CHARSET"],
		MaxIdleConns: MAXIDLECONNS,
		MaxOpenConns: MAXOPENCONNS,
		TimeZone:     "",
	}

	return response
}
