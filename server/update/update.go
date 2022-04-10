package update

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/Festivals-App/festivals-gateway/server/status"
)

func RunUpdate(repo string, updateScriptPath string) (string, error) {

	currentVersion := status.VersionString()
	if currentVersion == "development" {
		return "", errors.New("this is a development server please update manually")
	}

	newestVersion, err := LatestVersion(repo)
	if err != nil {
		return "", errors.New("failed to retrieve latest release version with error: " + err.Error())
	}

	if newestVersion == currentVersion {
		return "", errors.New("this server is already up-to-date")
	}

	cmd := exec.Command("/bin/bash", "-c", "sudo "+updateScriptPath+" &")
	err = cmd.Start()
	if err != nil {
		return "", errors.New("Failed to run update script with error: " + err.Error())
	}

	return newestVersion, nil
}

func LatestVersion(repository string) (string, error) {

	cmd := "curl --silent '" + repository + "' | sed -E 's/.*\"([^\"]+)\".*/\\1/' | xargs basename"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to execute command: %s with error: %s", cmd, err.Error())
		return "<Unknown>", errors.New(errorMessage)
	}
	return string(out), nil
}
