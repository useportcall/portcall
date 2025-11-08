package envx

import "github.com/joho/godotenv"

func Load() {
	var dotenv string

	if IsProd() {
		dotenv = ".env"
	} else {
		dotenv = ".env.example"
	}

	if err := godotenv.Load(dotenv); err != nil {
		panic("Error loading .env file")
	}
}
