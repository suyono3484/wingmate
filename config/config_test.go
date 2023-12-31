package config

import (
	"os"
	"path"
	"testing"

	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	serviceDir = "service"
)

var (
	configDir string
)

func setup(t *testing.T) {
	var err error
	if configDir, err = os.MkdirTemp("", "wingmate-*-test"); err != nil {
		t.Fatal("setup", err)
	}
	viper.Set(EnvConfigPath, configDir)
}

func tear(t *testing.T) {
	if err := os.RemoveAll(configDir); err != nil {
		t.Fatal("tear", err)
	}
}

func TestRead(t *testing.T) {

	type testEntry struct {
		name     string
		testFunc func(t *testing.T)
	}

	mkSvcDir := func(t *testing.T) {
		if err := os.MkdirAll(path.Join(configDir, serviceDir), 0755); err != nil {
			t.Fatal("create dir", err)
		}
	}

	touchFile := func(t *testing.T, name string, perm os.FileMode) {
		f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
		if err != nil {
			t.Fatal("create file", err)
		}
		_ = f.Close()
	}

	_ = wingmate.NewLog(os.Stderr)
	tests := []testEntry{
		{
			name: "positive",
			testFunc: func(t *testing.T) {
				mkSvcDir(t)
				touchFile(t, path.Join(configDir, serviceDir, "one.sh"), 0755)
				touchFile(t, path.Join(configDir, serviceDir, "two.sh"), 0755)

				cfg, err := Read()
				assert.Nil(t, err)
				assert.ElementsMatch(
					t,
					cfg.ServicePaths,
					[]string{
						path.Join(configDir, serviceDir, "one.sh"),
						path.Join(configDir, serviceDir, "two.sh"),
					},
				)
			},
		},
		{
			name: "with directory",
			testFunc: func(t *testing.T) {
				const subdir1 = "subdir1"
				mkSvcDir(t)
				assert.Nil(t, os.Mkdir(path.Join(configDir, serviceDir, subdir1), 0755))
				touchFile(t, path.Join(configDir, serviceDir, subdir1, "one.sh"), 0755)
				touchFile(t, path.Join(configDir, serviceDir, "two.sh"), 0755)
				cfg, err := Read()
				assert.Nil(t, err)
				assert.ElementsMatch(
					t,
					cfg.ServicePaths,
					[]string{
						path.Join(configDir, serviceDir, "two.sh"),
					},
				)
			},
		},
		{
			name: "wrong mode",
			testFunc: func(t *testing.T) {
				mkSvcDir(t)
				touchFile(t, path.Join(configDir, serviceDir, "one.sh"), 0755)
				touchFile(t, path.Join(configDir, serviceDir, "two.sh"), 0644)

				cfg, err := Read()
				assert.Nil(t, err)
				assert.ElementsMatch(
					t,
					cfg.ServicePaths,
					[]string{
						path.Join(configDir, serviceDir, "one.sh"),
					},
				)
			},
		},
		{
			name: "empty",
			testFunc: func(t *testing.T) {
				mkSvcDir(t)

				_, err := Read()
				assert.NotNil(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup(t)
			tt.testFunc(t)
			tear(t)
		})
	}
}
