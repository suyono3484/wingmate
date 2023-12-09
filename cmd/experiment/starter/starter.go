package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/spf13/viper"
)

const (
	// DummyPath = "/workspaces/wingmate/cmd/experiment/dummy/dummy"
	DummyPath    = "/usr/local/bin/wmdummy"
	EnvDummyPath = "DUMMY_PATH"
	EnvPrefix    = "WINGMATE"
)

func main() {
	var (
		stdout  io.ReadCloser
		stderr  io.ReadCloser
		wg      *sync.WaitGroup
		err     error
		exePath string
	)
	viper.SetEnvPrefix(EnvPrefix)
	viper.BindEnv(EnvDummyPath)
	viper.SetDefault(EnvDummyPath, DummyPath)

	exePath = viper.GetString(EnvDummyPath)

	cmd := exec.Command(exePath)

	if stdout, err = cmd.StdoutPipe(); err != nil {
		log.Panic(err)
	}

	if stderr, err = cmd.StderrPipe(); err != nil {
		log.Panic(err)
	}

	wg = &sync.WaitGroup{}
	wg.Add(2)
	go pulley(wg, stdout, "stdout")
	go pulley(wg, stderr, "stderr")

	if err = cmd.Start(); err != nil {
		log.Panic(err)
	}
	wg.Wait()

	if err = cmd.Wait(); err != nil {
		log.Printf("got error when Waiting for child process: %#v\n", err)
	}
}

func pulley(wg *sync.WaitGroup, src io.ReadCloser, srcName string) {
	defer wg.Done()

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		log.Printf("coming out from %s: %s\n", srcName, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("got error whean reading from %s: %#v\n", srcName, err)
	}

	log.Printf("closing %s...\n", srcName)
	_ = src.Close()
}
