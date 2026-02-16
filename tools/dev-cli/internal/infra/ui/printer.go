package ui

import (
	"fmt"
	"os"
)

const (
	blue   = "\033[34m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	bold   = "\033[1m"
	reset  = "\033[0m"
)

type Printer struct{}

func New() Printer { return Printer{} }

func (Printer) Plain(msg string, args ...any) { fmt.Printf(msg+"\n", args...) }
func (p Printer) Info(msg string, args ...any) {
	fmt.Printf(p.style("ℹ "+msg, blue)+"\n", args...)
}
func (p Printer) OK(msg string, args ...any) {
	fmt.Printf(p.style("✔ "+msg, green)+"\n", args...)
}
func (p Printer) Warn(msg string, args ...any) {
	fmt.Printf(p.style("⚠ "+msg, yellow)+"\n", args...)
}
func (p Printer) Fail(msg string, args ...any) {
	fmt.Printf(p.style("✖ "+msg, red)+"\n", args...)
}
func (p Printer) Section(title string) {
	fmt.Println(p.style("\n━━ "+title+" ━━", bold+blue))
}

func (Printer) style(text, color string) string {
	if os.Getenv("NO_COLOR") != "" {
		return text
	}
	return color + text + reset
}
