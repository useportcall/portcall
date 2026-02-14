package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetIconLogo(c *routerx.Context) {
	filename := c.Param("filename")

	// Reject directory traversal or suspicious filenames.
	if strings.ContainsAny(filename, "/\\") || filename != filepath.Base(filename) {
		c.String(http.StatusBadRequest, "Invalid filename")
		return
	}

	data, err := c.Store().GetFromIconLogoBucket(filename, c.Context)
	if err != nil {
		c.String(http.StatusNotFound, "Icon/Logo not found")
		return
	}

	contentType := http.DetectContentType(data)
	c.Data(http.StatusOK, contentType, data)
}
