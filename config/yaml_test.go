package config

import (
	"gitea.suyono.dev/suyono/wingmate"
	"os"
	"path"
	"testing"
)

const configName = "wingmate.yaml"

func TestYaml(t *testing.T) {
	type testEntry struct {
		name    string
		config  string
		wantErr bool
	}

	_ = wingmate.NewLog(os.Stderr)
	tests := []testEntry{
		{
			name:    "positive",
			config:  yamlTestCase0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup(t)
			defer tear(t)

			writeYaml(t, path.Join(configDir, configName), tt.config)

			cfg, err := Read()
			if tt.wantErr != (err != nil) {
				t.Fatalf("wantErr is %v but err is %+v", tt.wantErr, err)
			}
			t.Logf("cfg: %+v", cfg)
		})
	}
}

func writeYaml(t *testing.T, path, content string) {
	var (
		f   *os.File
		err error
	)

	if f, err = os.Create(path); err != nil {
		t.Fatal("create yaml file", err)
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err = f.Write([]byte(content)); err != nil {
		t.Fatal("write yaml file", err)
	}
}

const yamlTestCase0 = `version: "1"
service:
    one:
        command: ["command", "arg0", "arg1"]
        environ: ["ENV1=value1", "ENV2=valueX"]
        user: "user1"
        group: "999"
        working_dir: "/path/to/working"`
