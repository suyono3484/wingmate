package config

import (
	"fmt"
	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/viper"
)

const (
	ServiceConfigGroup = "service"
	CronConfigGroup    = "cron"
	ServiceKeyFormat   = "service.%s"
	CronKeyFormat      = "cron.%s"
)

func readConfigYaml(path, name, format string) ([]ServiceTask, []CronTask, error) {
	var (
		err         error
		nameMap     map[string]any
		itemName    string
		serviceTask ServiceTask
		cronTask    CronTask
		item        any
		services    []ServiceTask
		crones      []CronTask
	)

	viper.AddConfigPath(path)
	viper.SetConfigType(format)
	viper.SetConfigName(name)

	if err = viper.ReadInConfig(); err != nil {
		return nil, nil, fmt.Errorf("reading config in dir %s, file %s, format %s: %w", path, name, format, err)
	}

	services = make([]ServiceTask, 0)
	nameMap = viper.GetStringMap(ServiceConfigGroup)
	for itemName, item = range nameMap {
		serviceTask = ServiceTask{}
		if err = viper.UnmarshalKey(fmt.Sprintf(ServiceKeyFormat, itemName), &serviceTask); err != nil {
			wingmate.Log().Error().Msgf("failed to parse service %s: %+v | %+v", itemName, err, item)
			continue
		}
		serviceTask.Name = itemName
		services = append(services, serviceTask)
	}

	crones = make([]CronTask, 0)
	nameMap = viper.GetStringMap(CronConfigGroup)
	for itemName, item = range nameMap {
		cronTask = CronTask{}
		if err = viper.UnmarshalKey(fmt.Sprintf(CronKeyFormat, itemName), &cronTask); err != nil {
			wingmate.Log().Error().Msgf("failed to parse cron %s: %v | %v", itemName, err, item)
			continue
		}
		cronTask.Name = itemName
		crones = append(crones, cronTask)
	}

	return services, crones, nil
}
