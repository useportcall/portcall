package secrets

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runOut(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func askYesNo(prompt string, defaultYes bool) bool {
	suffix := " [y/N]: "
	if defaultYes {
		suffix = " [Y/n]: "
	}
	fmt.Print(prompt)
	if !strings.HasSuffix(prompt, ": ") {
		fmt.Print(suffix)
	}
	in, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	ans := strings.TrimSpace(strings.ToLower(in))
	if ans == "" {
		return defaultYes
	}
	return ans == "y" || ans == "yes"
}
