package config

import (
	"fmt"
	"os"
	"path"
	"testing"

	"gitea.suyono.dev/suyono/wingmate"
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
		{
			name:    "service only",
			config:  yamlTestCase1,
			wantErr: false,
		},
		{
			name:    "cron only",
			config:  yamlTestCase2,
			wantErr: false,
		},
		{
			name:    "invalid content - service",
			config:  yamlTestCase3,
			wantErr: false,
		},
	}

	for i, tc := range yamlBlobs {
		tests = append(tests, testEntry{
			name:    fmt.Sprintf("negative - %d", i),
			config:  tc,
			wantErr: true,
		})
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
        working_dir: "/path/to/working"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/5 * * * 2,3"`

const yamlTestCase1 = `version: "1"
service:
    one:
        command: ["command", "arg0", "arg1"]
        environ: ["ENV1=value1", "ENV2=valueX"]
        user: "user1"
        group: "999"
        working_dir: "/path/to/working"`

const yamlTestCase2 = `version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/5 * * * 2,3"`

const yamlTestCase3 = `version: "1"
service:
    one:
        command: 12345
        environ: ["ENV1=value1", "ENV2=valueX"]
        user: "user1"
        group: "999"
        working_dir: "/path/to/working"`

var yamlBlobs = []string{
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "a 13 3,5,7 * *"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/5 a 3,5,7 * *"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/5 13 a * *"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/5 13 3,5,7 a *"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/5 13 3,5,7 * a"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/x 13 3,5,7 * *"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "76 13 3,5,7 * *"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/75 13 3,5,7 * *"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/5 13 3,x,7 * *"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/5 13 3,5,67 * *"`,
	`version: "1"
cron:
    cron-one:
        command:
            - command-cron
            - arg0
            - arg1
        environ: ["ENV1=v1", "ENV2=var2"]
        user: "1001"
        group: "978"
        schedule: "*/5 13 * *"`,
}
