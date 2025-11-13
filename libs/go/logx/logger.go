package logx

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func Init() {
	env := os.Getenv("APP_ENV")

	// always log to stdout
	writers := []io.Writer{os.Stdout}

	if env == "" || env == "development" {
		// ensure logs dir exists
		if err := os.MkdirAll("logs", 0o755); err != nil {
			log.Printf("could not create logs dir: %v", err)
		} else {
			f, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				log.Printf("could not open dev log file: %v", err)
			} else {
				writers = append(writers, f)
			}
		}
	}

	gin.DefaultWriter = io.MultiWriter(writers...)
	// if you also want gin.DefaultErrorWriter:
	gin.DefaultErrorWriter = gin.DefaultWriter
}
