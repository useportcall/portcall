package middleware

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/useportcall/portcall/libs/go/routerx"
)

func StaticFile(staticDir string) routerx.HandlerFunc {
	return func(c *routerx.Context) {
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

		// Heuristic: only SPA-fallback when the client likely wants HTML
		accept := c.GetHeader("Accept")
		wantsHTML := strings.Contains(accept, "text/html") || !strings.Contains(path.Base(reqPath), ".")

		if wantsHTML {
			http.ServeFile(c.Writer, c.Request, filepath.Join(staticDir, "index.html"))
			c.Abort()
			return
		}

		// Otherwise, it's probably a missing asset → 404
		c.AbortWithStatus(http.StatusNotFound)
	}
}
