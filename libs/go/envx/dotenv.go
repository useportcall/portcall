package envx

import (
	"log"

	"github.com/joho/godotenv"
)

func Load() {
	var dotenv string

	if IsProd() {
		dotenv = ".env"
	} else {
		dotenv = ".env.example"
	}

	if err := godotenv.Load(dotenv); err != nil {
		log.Printf("could not load %s file: %v", dotenv, err)
	}
}
