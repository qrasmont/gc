package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

func runGitBranch() ([]byte, error) {
	shell := os.Getenv("SHELL")

	out, err := exec.Command(shell, "-c", "git branch --no-color").Output()
	if err != nil {
		return nil, err
	}
	return out, err
}

func getBranchList(raw_b []byte) []list.Item {
	var items []list.Item
	raw_s := string(raw_b)
	branches := strings.Split(strings.ReplaceAll(raw_s, "\r\n", "\n"), "\n")
	// Remove last element which is an empty line
	branches = branches[:len(branches)-1]

	for i, branch := range branches {
		branches[i] = strings.TrimLeft(branch, "*")
		branches[i] = strings.TrimLeft(branches[i], " ")
        items = append(items, item{name: branches[i]})
	}

	return items
}

func GitGetBranches() ([]list.Item, error) {
	raw_b, err := runGitBranch()
	if err != nil {
		return nil, err
	}

	list := getBranchList(raw_b)
	return list, nil
}

func GitDelete(selection []string) error {
	shell := os.Getenv("SHELL")
	branches := fmt.Sprintf("%s", strings.Join(selection[:], " "))

	err := exec.Command(shell, "-c", "git branch -D "+branches).Run()
	if err != nil {
		return err
	}
	return nil
}
