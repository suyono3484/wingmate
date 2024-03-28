package config

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func getVersion(binPath string) (string, error) {
	var (
		outBytes []byte
		err      error
		output   string
		stdout   io.ReadCloser
		n        int
	)
	cmd := exec.Command(binPath, "version")
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return "", fmt.Errorf("setting up stdout reader: %w", err)
	}

	if err = cmd.Start(); err != nil {
		return "", fmt.Errorf("starting process: %w", err)
	}

	outBytes = make([]byte, 1024)
	if n, err = stdout.Read(outBytes); err != nil {
		return "", fmt.Errorf("reading stdout: %w", err)
	}

	_ = cmd.Wait()

	output = string(outBytes[:n])
	output = strings.TrimRight(output, versionTrimRightCutSet)
	return output, nil
}
