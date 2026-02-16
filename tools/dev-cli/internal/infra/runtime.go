package infra

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	rootDir      string
	rootResolver = findRootDir
)

func SetRootResolver(fn func() (string, error)) {
	if fn != nil {
		rootResolver = fn
	}
}

func ensureRootDir() error {
	if strings.TrimSpace(rootDir) != "" {
		return nil
	}
	dir, err := rootResolver()
	if err != nil {
		return err
	}
	rootDir = dir
	return nil
}

func findRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("could not find go.work in any parent directory")
}

func runCmd(name string, args ...string) error {
	if err := ensureRootDir(); err != nil {
		return err
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = rootDir
	return cmd.Run()
}

func runCmdWithEnv(env map[string]string, name string, args ...string) error {
	if err := ensureRootDir(); err != nil {
		return err
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = rootDir
	cmd.Env = mergeCmdEnv(env)
	return cmd.Run()
}

func runCmdOut(name string, args ...string) (string, error) {
	if err := ensureRootDir(); err != nil {
		return "", err
	}
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = rootDir
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runCmdOutWithEnv(env map[string]string, name string, args ...string) (string, error) {
	if err := ensureRootDir(); err != nil {
		return "", err
	}
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = rootDir
	cmd.Env = mergeCmdEnv(env)
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func mergeCmdEnv(override map[string]string) []string {
	env := os.Environ()
	for key, value := range override {
		env = append(env, key+"="+value)
	}
	return env
}

func runShell(script string) error {
	return runCmd("bash", "-lc", script)
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.TrimSpace(line)
}

func switchCluster(cluster string) error {
	if err := ensureRootDir(); err != nil {
		return err
	}
	return runCmd("kubectl", "config", "use-context", ResolveClusterContext(rootDir, cluster))
}
