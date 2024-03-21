package main

import (
	"os"
	"testing"
)

func TestEntry_configPathEnv(t *testing.T) {
	_ = os.Setenv("WINGMATE_CONFIG_PATH", "/Volumes/Source/go/src/gitea.suyono.dev/suyono/wingmate/docker/bookworm/etc/wingmate")
	defer func() {
		_ = os.Unsetenv("WINGMATE_CONFIG_PATH")
	}()
	main()
}

func TestEntry_configPathPFlag(t *testing.T) {
	os.Args = []string{"wingmate", "--config", "/Volumes/Source/go/src/gitea.suyono.dev/suyono/wingmate/docker/bookworm/etc/wingmate"}
	main()
}

func TestEntry(t *testing.T) {
	main()
}
