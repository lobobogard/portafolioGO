package env

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Env() map[string]string {
	var Env = make(map[string]string)

	err := godotenv.Load(filepath.Join("/root/", ".env"))

	if err != nil {
		log.Fatal("Error loading .env emv.go", err)
	}

	envPath, _ := filepath.Abs("/root/.env")
	Env, err = godotenv.Read(envPath)
	if err != nil {
		log.Fatal("error reading .env env.go ", err)
	}

	return Env

}
