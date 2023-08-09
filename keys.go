package main

import (
	"fmt"
	"regexp"
)

// getKeys returns the list of keys from the current running local node
func getKeys() ([]string, error) {
	out, err := executeShellCommand([]string{"keys", "list"}, evmosdHome, "", false, false)
	if err != nil {
		return nil, err
	}

	return parseKeysFromOut(out)
}

func parseKeysFromOut(out string) ([]string, error) {
	// Define the regular expression pattern
	pattern := `\s+name:\s*(\w+)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(out, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no keys found in output")
	}

	var keys []string
	for _, match := range matches {
		keys = append(keys, match[1])
	}

	return keys, nil
}
