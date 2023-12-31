package config

import (
	"gitea.suyono.dev/suyono/wingmate"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

const (
	crontabFileName = "crontab"
)

func TestCrontab(t *testing.T) {
	type testEntry struct {
		name    string
		crontab string
		wantErr bool
	}

	_ = wingmate.NewLog(os.Stderr)
	tests := []testEntry{
		{
			name:    "positive",
			crontab: crontabTestCase0,
			wantErr: false,
		},
		{
			name:    "with comment",
			crontab: crontabTestCase1,
			wantErr: false,
		},
		{
			name:    "various values",
			crontab: crontabTestCase2,
			wantErr: false,
		},
		{
			name:    "failed to parse",
			crontab: crontabTestCase3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup(t)
			defer tear(t)

			writeCrontab(t, tt.crontab)

			cfg, err := Read()
			if tt.wantErr != (err != nil) {
				t.Fatalf("wantErr is %v but err is %+v", tt.wantErr, err)
			}

			t.Logf("cfg: %+v", cfg)
			for _, c := range cfg.Cron {
				t.Logf("%+v", c)
			}
		})
	}
}

func writeCrontab(t *testing.T, content string) {
	var (
		f   *os.File
		err error
	)

	if f, err = os.Create(filepath.Join(configDir, crontabFileName)); err != nil {
		t.Fatal("create crontab file", err)
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err = f.Write([]byte(content)); err != nil {
		t.Fatal("writing crontab file", err)
	}
}

const crontabTestCase0 = `* * * * *  /path/to/executable`
const crontabTestCase1 = `# this is a comment
  ## comment with space
* * * * *  /path/to/executable
* * * * *  /path/to/executable  # comment as a suffix
`

const crontabTestCase2 = `# first comment
*/5 13 3,5,7 * *  /path/to/executable`

const crontabTestCase3 = `a 13 3,5,7 * *  /path/to/executable
*/5 a 3,5,7 * *  /path/to/executable
*/5 13 a * *  /path/to/executable
*/5 13 3,5,7 a *  /path/to/executable
*/5 13 3,5,7 * a  /path/to/executable
*/x 13 3,5,7 * a  /path/to/executable
76 13 3,5,7 * a  /path/to/executable
*/75 13 3,5,7 * a  /path/to/executable
*/5 13 3,x,7 * a  /path/to/executable
*/5 13 3,5,67 * a  /path/to/executable
*/5 13 * *  /path/to/executable
*/5 13 3,5,7 * *  /path/to/executable`

func TestSpecExact(t *testing.T) {
	var val uint8 = 45
	s := SpecExact{
		value: val,
	}

	assert.Equal(t, val, s.Value())
}

func TestSpecMulti(t *testing.T) {
	val := []uint8{3, 5, 7, 15}
	s := SpecMultiOccurrence{
		values: val,
	}

	assert.ElementsMatch(t, val, s.Values())
}

func TestInvalidField(t *testing.T) {
	c := &Cron{}
	assert.NotNil(t, c.setField(cronField(99), "x"))
}
