package envx

import (
	"log"

	"github.com/joho/godotenv"
)

func Load() {
	files := []string{".envs", ".env"}
	if !IsProd() {
		files = append(files, ".env.example")
	}

	for _, file := range files {
		if err := godotenv.Load(file); err != nil {
			log.Printf("could not load %s file: %v", file, err)
		}
	}
}
