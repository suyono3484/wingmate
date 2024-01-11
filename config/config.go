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
	EnvPrefix                = "WINGMATE"
	EnvConfigPath            = "CONFIG_PATH"
	DefaultConfigPath        = "/etc/wingmate"
	ServiceDirName           = "service"
	CrontabFileName          = "crontab"
	WingmateConfigFileName   = "wingmate"
	WingmateConfigFileFormat = "yaml"
)

type Config struct {
	ServiceV0 []string
	CronV0    []*Cron
	Service   []ServiceTask
	Cron      []CronTask
}

type Task struct {
	Command    []string `mapstructure:"command"`
	Environ    []string `mapstructure:"environ"`
	Setsid     bool     `mapstructure:"setsid"`
	User       string   `mapstructure:"user"`
	Group      string   `mapstructure:"group"`
	WorkingDir string   `mapstructure:"working_dir"`
}

type ServiceTask struct {
	Task        `mapstructure:",squash"`
	Name        string `mapstructure:"-"`
	Background  bool   `mapstructure:"background"`
	PidFile     string `mapstructure:"pidfile"`
	StartSecs   uint   `mapstructure:"startsecs"`
	AutoStart   bool   `mapstructure:"autostart"`
	AutoRestart bool   `mapstructure:"autorestart"`
}

type CronTask struct {
	CronSchedule `mapstructure:"-"`
	Task         `mapstructure:",squash"`
	Name         string `mapstructure:"-"`
	Schedule     string `mapstructure:"schedule"`
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
		dirent                  []os.DirEntry
		err                     error
		svcdir                  string
		serviceAvailable        bool
		cronAvailable           bool
		wingmateConfigAvailable bool
		cron                    []*Cron
		crontabfile             string
		services                []ServiceTask
		crones                  []CronTask
	)

	serviceAvailable = false
	cronAvailable = false
	outConfig := &Config{
		ServiceV0: make([]string, 0),
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
					outConfig.ServiceV0 = append(outConfig.ServiceV0, svcPath)
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
		outConfig.CronV0 = cron
		cronAvailable = true
	}
	if err != nil {
		wingmate.Log().Error().Msgf("encounter error when reading crontab %s: %+v", crontabfile, err)
	}

	wingmateConfigAvailable = false
	if services, crones, err = readConfigYaml(configPath, WingmateConfigFileName, WingmateConfigFileFormat); err != nil {
		wingmate.Log().Error().Msgf("encounter error when reading wingmate config file in %s/%s: %+v", configPath, WingmateConfigFileName, err)
	}
	if len(services) > 0 {
		outConfig.Service = services
		wingmateConfigAvailable = true
	}
	if len(crones) > 0 {
		outConfig.Cron = crones
		wingmateConfigAvailable = true
	}

	if !serviceAvailable && !cronAvailable && !wingmateConfigAvailable {
		return nil, errors.New("no config found")
	}

	return outConfig, nil
}
