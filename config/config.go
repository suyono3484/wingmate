package config

import (
	"errors"
	"os"
	"path/filepath"

	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/viper"
)

const (
	EnvPrefix         = "WINGMATE"
	EnvConfigPath     = "CONFIG_PATH"
	DefaultConfigPath = "/etc/wingmate"
	ServiceDirName    = "service"
	CrontabFileName   = "crontab"
)

type Config struct {
	servicePaths []string
}

type cron struct{}

func Read() (*Config, error) {
	viper.SetEnvPrefix(EnvPrefix)
	viper.BindEnv(EnvConfigPath)
	viper.SetDefault(EnvConfigPath, DefaultConfigPath)

	var (
		dirent           []os.DirEntry
		err              error
		svcdir           string
		serviceAvailable bool
		cronAvailable    bool
		cron             []cron
		crontabfile      string
	)

	serviceAvailable = false
	cronAvailable = false
	outConfig := &Config{
		servicePaths: make([]string, 0),
	}
	configPath := viper.GetString(EnvConfigPath)
	svcdir = filepath.Join(configPath, ServiceDirName)
	dirent, err = os.ReadDir(svcdir)
	if len(dirent) > 0 {
		for _, d := range dirent {
			if d.Type().IsRegular() {
				serviceAvailable = true
				outConfig.servicePaths = append(outConfig.servicePaths, filepath.Join(svcdir, d.Name()))
			}
		}
	}
	if err != nil {
		wingmate.Log().Error().Msgf("encounter error when reading service directory %s: %#v", svcdir, err)
	}

	crontabfile = filepath.Join(configPath, CrontabFileName)
	cron, err = readCrontab(crontabfile)
	if len(cron) > 0 {
		cronAvailable = true
	}
	if err != nil {
		wingmate.Log().Error().Msgf("encounter error when reading crontab %s: %#v", crontabfile, err)
	}

	if !serviceAvailable && !cronAvailable {
		return nil, errors.New("no config found")
	}

	return outConfig, nil
}
