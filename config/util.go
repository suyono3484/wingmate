package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os/exec"
	"strings"
)

type pathResult struct {
	path string
	err  error
}

type FindUtils struct {
	exec     chan pathResult
	pidProxy chan pathResult
}

const (
	PidProxyBinaryName = "wmpidproxy"
	ExecBinaryName     = "wmexec"
)

func findExec(currentVersion string, path string, resultChan chan<- pathResult) {
	var (
		result      pathResult
		execVersion string
	)
	defer close(resultChan)

	if len(path) > 0 {
		result.path = path
		execVersion, result.err = getVersion(result.path)
	} else {
		result.path, result.err = exec.LookPath(ExecBinaryName)
		if result.err != nil {
			resultChan <- result
			return
		}
		execVersion, result.err = getVersion(result.path)
	}

	if result.err == nil {
		if execVersion != currentVersion {
			result.err = fmt.Errorf("incompatible version: wingmate %s and wmexec %s", currentVersion, execVersion)
		}
	}

	resultChan <- result
}

func findPidProxy(currentVersion string, path string, resultChan chan<- pathResult) {
	var (
		result          pathResult
		pidProxyVersion string
	)
	defer close(resultChan)

	if len(path) > 0 {
		result.path = path
		pidProxyVersion, result.err = getVersion(result.path)
	} else {
		result.path, result.err = exec.LookPath(PidProxyBinaryName)
		if result.err != nil {
			resultChan <- result
			return
		}
		pidProxyVersion, result.err = getVersion(result.path)
	}

	if result.err == nil {
		if pidProxyVersion != currentVersion {
			result.err = fmt.Errorf("incompatible version: wingmate %s and wmpidproxy %s", currentVersion, pidProxyVersion)
		}
	}

	resultChan <- result
}

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

func startFindUtils() *FindUtils {
	result := &FindUtils{
		exec:     make(chan pathResult),
		pidProxy: make(chan pathResult),
	}

	var (
		pidProxyPath string
		execPath     string
	)

	currentVersion := viper.GetString(WingmateVersion)

	if viper.IsSet(PidProxyPathConfig) {
		pidProxyPath = viper.GetString(PidProxyPathConfig)
	}

	if viper.IsSet(ExecPathConfig) {
		execPath = viper.GetString(ExecPathConfig)
	}

	go findPidProxy(currentVersion, pidProxyPath, result.pidProxy)
	go findExec(currentVersion, execPath, result.exec)
	return result
}

//func (f *FindUtils) Get
