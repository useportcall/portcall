package routerx

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func StaticFileMiddleware(staticDir string) HandlerFunc {
	return func(c *Context) {
		// Let API and non-GET requests through untouched
		if strings.HasPrefix(c.Request.URL.Path, "/api") || c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Clean + join; prevent path traversal
		reqPath := c.Request.URL.Path
		full := filepath.Join(staticDir, filepath.Clean(reqPath))
		if rel, err := filepath.Rel(staticDir, full); err != nil || strings.HasPrefix(rel, "..") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// If a real asset exists, serve it
		if fi, err := os.Stat(full); err == nil && !fi.IsDir() {
			http.ServeFile(c.Writer, c.Request, full)
			c.Abort()
			return
		}

		// Next.js static export: try {path}.html for clean URLs (e.g. /success → success.html)
		if !strings.Contains(path.Base(reqPath), ".") {
			htmlPath := full + ".html"
			if fi, err := os.Stat(htmlPath); err == nil && !fi.IsDir() {
				http.ServeFile(c.Writer, c.Request, htmlPath)
				c.Abort()
				return
			}
		}

		// Heuristic: only SPA-fallback when the client likely wants HTML
		accept := c.GetHeader("Accept")
		wantsHTML := strings.Contains(accept, "text/html") || !strings.Contains(path.Base(reqPath), ".")

		if wantsHTML {
			// Inject environment variables into index.html
			indexPath := filepath.Join(staticDir, "index.html")
			htmlContent, err := os.ReadFile(indexPath)
			if err != nil {
				http.ServeFile(c.Writer, c.Request, indexPath)
				c.Abort()
				return
			}

			// Inject environment variables as a script tag
			envScript := fmt.Sprintf(`<script>window.__ENV__ = {
  QUOTE_APP_URL: "%s",
  KEYCLOAK_URL: "%s",
  API_URL: "%s",
  FILE_API_URL: "%s"
};</script>`,
				getEnvOrDefault("QUOTE_APP_URL", "http://localhost:8010"),
				getEnvOrDefault("KEYCLOAK_API_URL", "http://localhost:8080"),
				getEnvOrDefault("API_URL", "http://localhost:8080"),
				getEnvOrDefault("FILE_API_URL", "http://localhost:8085"),
			)

			// Inject E2E mode variables when present
			if e2e := os.Getenv("E2E_MODE"); e2e != "" {
				envScript += fmt.Sprintf(`<script>Object.assign(window.__ENV__,{E2E_MODE:true,E2E_ADMIN_KEY:"%s",E2E_APP_ID:"%s"});</script>`,
					os.Getenv("ADMIN_API_KEY"),
					os.Getenv("E2E_APP_ID"),
				)
			}

			// Insert before </head>
			modifiedHTML := bytes.Replace(htmlContent, []byte("</head>"), []byte(envScript+"\n  </head>"), 1)

			c.Data(http.StatusOK, "text/html; charset=utf-8", modifiedHTML)
			c.Abort()
			return
		}

		// Otherwise, it's probably a missing asset → 404
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
