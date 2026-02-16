package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type choiceItem struct {
	Key   string
	Label string
	Value string
}

var menuReader = bufio.NewReader(os.Stdin)
var errBack = errors.New("back")

func promptChoice(title string, items []choiceItem) (string, error) {
	fmt.Printf("\n== %s ==\n", title)
	for _, item := range items {
		fmt.Printf("  %s) %s\n", item.Key, item.Label)
	}
	for {
		in, err := promptLine("Choose: ")
		if err != nil {
			return "", err
		}
		for _, item := range items {
			if strings.EqualFold(in, item.Key) {
				return item.Value, nil
			}
		}
		fmt.Println("Invalid choice. Try again.")
	}
}

func promptLine(prompt string) (string, error) {
	fmt.Print(prompt)
	line, err := menuReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func promptWithDefault(label, fallback string) (string, error) {
	in, err := promptLine(fmt.Sprintf("%s [%s]: ", label, fallback))
	if err != nil {
		return "", err
	}
	switch strings.ToLower(strings.TrimSpace(in)) {
	case "back", "cancel":
		return "", errBack
	}
	if strings.TrimSpace(in) == "" {
		return fallback, nil
	}
	return in, nil
}
