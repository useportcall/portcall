package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/useportcall/portcall/apps/api/internal/router"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func main() {
	envx.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	db, err := dbx.New()
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}
	crypto, err := cryptox.New()
	if err != nil {
		log.Fatalf("failed to init crypto: %v", err)
	}
	q, err := qx.New()
	if err != nil {
		log.Fatalf("failed to init queue: %v", err)
	}

	r := router.Init(db, crypto, q)

	// Serve OpenAPI spec at /openapi.yaml (public, no auth required)
	r.GET("/openapi.yaml", func(c *routerx.Context) {
		c.Writer.Header().Set("Content-Type", "application/x-yaml")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(c.Writer, c.Request, "./openapi.yaml")
	})

	// Handle /docs and /docs/ - NoRoute handler as fallback
	r.NoRoute(func(c *routerx.Context) {
		path := c.Request.URL.Path

		if path == "/docs" {
			http.Redirect(c.Writer, c.Request, "/docs/", http.StatusMovedPermanently)
			return
		}

		if strings.HasPrefix(path, "/docs/") || path == "/docs" {
			reqPath := strings.TrimPrefix(path, "/docs")
			if reqPath == "" || reqPath == "/" {
				reqPath = "/index.html"
			}

			fullPath := filepath.Join("./frontend/dist", filepath.Clean(reqPath))

			if _, err := os.Stat(fullPath); err == nil {
				http.ServeFile(c.Writer, c.Request, fullPath)
				return
			}

			http.ServeFile(c.Writer, c.Request, "./frontend/dist/index.html")
			return
		}

		c.JSON(404, map[string]any{"error": "Not found"})
	})

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
