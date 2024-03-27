package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

const (
	EnvPrefix                = "WINGMATE"
	ConfigPath               = "config_path"
	DefaultConfigPath        = "/etc/wingmate"
	ServiceDirName           = "service"
	CrontabFileName          = "crontab"
	WingmateConfigFileName   = "wingmate"
	WingmateConfigFileFormat = "yaml"
	WingmateVersion          = "APP_VERSION"
	PidProxyPathConfig       = "pidproxy_path"
	PidProxyPathDefault      = "wmpidproxy"
	ExecPathConfig           = "exec_path"
	ExecPathDefault          = "wmexec"
	versionTrimRightCutSet   = "\r\n "
	VersionFlag              = "version"
	WMPidProxyPathFlag       = "pid-proxy"
	WMExecPathFlag           = "exec"
	ConfigPathFlag           = "config"
	VersionCheckKey          = "check-version"
)

type Config struct {
	ServiceV0 []string
	CronV0    []*Cron
	Service   []ServiceTask
	Cron      []CronTask
	viperMtx  *sync.Mutex
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

func SetVersion(version string) {
	version = strings.TrimRight(version, versionTrimRightCutSet)
	viper.Set(WingmateVersion, version)
	wingmate.Log().Info().Msgf("starting wingmate version %s", version)
}

func Read() (*Config, error) {
	viper.SetEnvPrefix(EnvPrefix)
	_ = viper.BindEnv(ConfigPath)
	_ = viper.BindEnv(PidProxyPathConfig)
	_ = viper.BindEnv(ExecPathConfig)
	viper.SetDefault(ConfigPath, DefaultConfigPath)
	viper.SetDefault(PidProxyPathConfig, PidProxyPathDefault)
	viper.SetDefault(ExecPathConfig, ExecPathDefault)

	ParseFlags()

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
		viperMtx:  &sync.Mutex{},
		ServiceV0: make([]string, 0),
	}
	configPath := viper.GetString(ConfigPath)
	svcdir = filepath.Join(configPath, ServiceDirName)
	dirent, err = os.ReadDir(svcdir)
	if err != nil {
		wingmate.Log().Error().Msgf("encounter error when reading service directory %s: %+v", svcdir, err)
	}
	if len(dirent) > 0 {
		for _, d := range dirent {
			if d.Type().IsRegular() {
				svcPath := filepath.Join(svcdir, d.Name())
				if err = unix.Access(svcPath, unix.X_OK); err == nil {
					serviceAvailable = true
					outConfig.ServiceV0 = append(outConfig.ServiceV0, svcPath)
				} else {
					wingmate.Log().Error().Msgf("checking executable access for %s: %+v", svcPath, err)
				}
			}
		}
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

func (c *Config) GetAppVersion() string {
	c.viperMtx.Lock()
	defer c.viperMtx.Unlock()

	return viper.GetString(WingmateVersion)
}

func (c *Config) WMPidProxyPath() string {
	c.viperMtx.Lock()
	defer c.viperMtx.Unlock()

	return viper.GetString(PidProxyPathConfig)
}

func (c *Config) WMPidProxyCheckVersion() error {
	var (
		binVersion string
		appVersion string
		err        error
	)

	if binVersion, err = getVersion(c.WMPidProxyPath()); err != nil {
		return fmt.Errorf("get wmpidproxy version: %w", err)
	}

	appVersion = c.GetAppVersion()
	if appVersion != binVersion {
		return fmt.Errorf("wmpidproxy version mismatch")
	}
	return nil
}

func (c *Config) WMExecPath() string {
	c.viperMtx.Lock()
	defer c.viperMtx.Unlock()

	return viper.GetString(ExecPathConfig)
}

func (c *Config) WMExecCheckVersion() error {
	var (
		binVersion string
		appVersion string
		err        error
	)

	if binVersion, err = getVersion(c.WMExecPath()); err != nil {
		return fmt.Errorf("get wmexec version: %w", err)
	}

	appVersion = c.GetAppVersion()
	if appVersion != binVersion {
		return fmt.Errorf("wmexec version mismatch")
	}

	return nil
}
