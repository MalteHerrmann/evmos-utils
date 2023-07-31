package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// executeShellCommand executes a shell command and returns the output and error.
func executeShellCommand(command []string, home string, sender string, defaults bool) (string, error) {
	fullCommand := command
	if home != "" {
		fullCommand = append(fullCommand, "--home", home)
	}
	if sender != "" {
		fullCommand = append(fullCommand, "--from", sender)
	}
	if defaults {
		fullCommand = append(fullCommand, defaultFlags...)
	}

	cmd := exec.Command("evmosd", fullCommand...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
	}
	return string(output), err
}

// getCurrentHeight returns the current block height of the node.
func getCurrentHeight() (int, error) {
	output, err := executeShellCommand([]string{"q", "block", "--node", "http://localhost:26657"}, evmosdHome, "", false)
	if err != nil {
		return 0, fmt.Errorf("error executing command: %w", err)
	}

	heightPattern := regexp.MustCompile(`"last_commit":{"height":"(\d+)"`)
	match := heightPattern.FindStringSubmatch(output)
	if len(match) < 2 {
		return 0, fmt.Errorf("did not find block height in output: \n%s", output)
	}

	height, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("error converting height to integer: %w", err)
	}

	return height, nil
}
