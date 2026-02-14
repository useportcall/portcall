package handlers

import "path/filepath"

// TemplateDir is the base directory for HTML templates.
// Override before starting the server when templates are not in ./templates.
var TemplateDir = "templates"

// tmplPaths returns absolute paths for the given template file names.
func tmplPaths(names ...string) []string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = filepath.Join(TemplateDir, n)
	}
	return out
}
