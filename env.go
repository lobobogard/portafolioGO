package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

var Env = make(map[string]string)

func addEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Env, err = godotenv.Read()

	if err != nil {
		log.Fatal("Error reading .env file")
	}
	// fmt.Println(Env["S3_BUCKET"])
	fmt.Println("File env charge correct")

}
