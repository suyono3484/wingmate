package config

import (
	"os/exec"
	"strings"
)

func getVersion(binPath string) (string, error) {
	var (
		outBytes []byte
		err      error
		output   string
	)
	cmd := exec.Command(binPath, "version")
	outBytes, err = cmd.Output()
	if err != nil {
		return "", err
	}

	output = string(outBytes)
	output = strings.TrimRight(output, versionTrimRightCutSet)
	return output, nil
}
