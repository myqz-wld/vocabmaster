package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vocabmaster/vocabmaster/src/internal/model"
)

var scanner = bufio.NewScanner(os.Stdin)

func WaitForEnter(hint string) {
	if hint == "" {
		hint = "按 Enter 继续..."
	}
	fmt.Printf("  %s%s%s ", colorGray, hint, colorReset)
	if !scanner.Scan() {
		os.Exit(0)
	}
}

func ReadGrade() model.Grade {
	for {
		DisplayGradeOptions()
		fmt.Print("  > ")
		if !scanner.Scan() {
			os.Exit(0)
		}
		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "0":
			return model.GradeAgain
		case "1":
			return model.GradeHard
		case "2":
			return model.GradeGood
		case "3":
			return model.GradeEasy
		default:
			fmt.Printf("  %s请输入 0-3%s\n", colorRed, colorReset)
		}
	}
}

func ReadConfirmation(prompt string) bool {
	fmt.Printf("  %s (y/N) ", prompt)
	if !scanner.Scan() {
		return false
	}
	input := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return input == "y" || input == "yes"
}
