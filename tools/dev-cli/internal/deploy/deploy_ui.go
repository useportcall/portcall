package deploy

import "fmt"

const (
	cRed    = "\033[0;31m"
	cGreen  = "\033[0;32m"
	cYellow = "\033[1;33m"
	cBlue   = "\033[0;34m"
	cCyan   = "\033[0;36m"
	cNC     = "\033[0m"
)

func info(msg string, args ...any)  { fmt.Printf(cBlue+msg+cNC+"\n", args...) }
func ok(msg string, args ...any)    { fmt.Printf(cGreen+msg+cNC+"\n", args...) }
func warn(msg string, args ...any)  { fmt.Printf(cYellow+msg+cNC+"\n", args...) }
func fail(msg string, args ...any)  { fmt.Printf(cRed+msg+cNC+"\n", args...) }
func plain(msg string, args ...any) { fmt.Printf(msg+"\n", args...) }
