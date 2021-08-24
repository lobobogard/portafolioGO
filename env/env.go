package env

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func Env() map[string]string {
	var Env = make(map[string]string)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Env, err = godotenv.Read()

	if err != nil {
		log.Fatal("Error reading .env file")
	}
	// fmt.Println(Env["database"])
	fmt.Println("File env charge correct")

	return Env

}
