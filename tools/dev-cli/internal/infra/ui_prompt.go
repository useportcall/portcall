package infra

import (
	"os"
	"strings"

	infraUI "github.com/useportcall/portcall/tools/dev-cli/internal/infra/ui"
)

var printer = infraUI.New()

func plain(msg string, args ...any) { printer.Plain(msg, args...) }
func info(msg string, args ...any)  { printer.Info(msg, args...) }
func ok(msg string, args ...any)    { printer.OK(msg, args...) }
func warn(msg string, args ...any)  { printer.Warn(msg, args...) }
func fail(msg string, args ...any)  { printer.Fail(msg, args...) }
func section(title string)          { printer.Section(title) }

func askYesNo(prompt string, defaultYes bool) bool {
	if !isInteractiveSession() {
		return defaultYes
	}
	for {
		answer := strings.ToLower(strings.TrimSpace(readInput(prompt)))
		if answer == "" {
			return defaultYes
		}
		if answer == "y" || answer == "yes" {
			return true
		}
		if answer == "n" || answer == "no" {
			return false
		}
		warn("Please answer y or n")
	}
}

func askText(prompt string, fallback string) string {
	if !isInteractiveSession() {
		return fallback
	}
	text := strings.TrimSpace(readInput(prompt))
	if text == "" {
		return fallback
	}
	return text
}

func isInteractiveSession() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
