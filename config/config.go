package config

import (
	"errors"
	"os"
	"path/filepath"

	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

const (
	EnvPrefix         = "WINGMATE"
	EnvConfigPath     = "CONFIG_PATH"
	DefaultConfigPath = "/etc/wingmate"
	ServiceDirName    = "service"
	CrontabFileName   = "crontab"
)

type Config struct {
	ServicePaths []string
	Cron         []*Cron
}

type Task struct {
	Command    []string `yaml:"command"`
	Environ    []string `yaml:"environ"`
	Setsid     bool     `yaml:"setsid"`
	User       string   `yaml:"user"`
	Group      string   `yaml:"group"`
	Background bool     `yaml:"background"`
	WorkingDir string   `yaml:"working_dir"`
}

type ServiceTask struct {
	Task        `yaml:",inline"`
	AutoStart   bool `yaml:"autostart"`
	AutoRestart bool `yaml:"autorestart"`
}

type CronTask struct {
	CronSchedule `yaml:"-"`
	Task         `yaml:",inline"`
	Schedule     string `yaml:"schedule"`
}

type CronSchedule struct {
	Minute CronTimeSpec
	Hour   CronTimeSpec
	DoM    CronTimeSpec
	Month  CronTimeSpec
	DoW    CronTimeSpec
}

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
		cron             []*Cron
		crontabfile      string
	)

	serviceAvailable = false
	cronAvailable = false
	outConfig := &Config{
		ServicePaths: make([]string, 0),
	}
	configPath := viper.GetString(EnvConfigPath)
	svcdir = filepath.Join(configPath, ServiceDirName)
	dirent, err = os.ReadDir(svcdir)
	if len(dirent) > 0 {
		for _, d := range dirent {
			if d.Type().IsRegular() {
				svcPath := filepath.Join(svcdir, d.Name())
				if err = unix.Access(svcPath, unix.X_OK); err == nil {
					serviceAvailable = true
					outConfig.ServicePaths = append(outConfig.ServicePaths, svcPath)
				}
			}
		}
	}
	if err != nil {
		wingmate.Log().Error().Msgf("encounter error when reading service directory %s: %+v", svcdir, err)
	}

	crontabfile = filepath.Join(configPath, CrontabFileName)
	cron, err = readCrontab(crontabfile)
	if len(cron) > 0 {
		outConfig.Cron = cron
		cronAvailable = true
	}
	if err != nil {
		wingmate.Log().Error().Msgf("encounter error when reading crontab %s: %+v", crontabfile, err)
	}

	if !serviceAvailable && !cronAvailable {
		return nil, errors.New("no config found")
	}

	return outConfig, nil
}
