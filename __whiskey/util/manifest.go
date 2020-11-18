package util

import (
	"log"
	"os/exec"
	"strings"
)

var cachedVersion string

// Version gets version data about whiskey.
func Version() string {
	if cachedVersion == "" {
		output, err := getLatestTag()
		if err != nil {
			if eErr, ok := err.(*exec.ExitError); ok {
				errOut := string(eErr.Stderr)
				if strings.Contains(errOut, "No names found") {
					cachedVersion = "development"
					return cachedVersion
				}
			}

			log.Fatalln("Failed to get latest git tag:", err)
		}

		cachedVersion = output
	}

	return cachedVersion
}

func getLatestTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--long", "`git rev-list --tags --max-count=1`")
	output, err := cmd.Output()
	return string(output), err
}
