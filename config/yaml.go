package config

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gitea.suyono.dev/suyono/wingmate"
	"github.com/spf13/viper"
)

const (
	CrontabScheduleRegexPattern = `^\s*(?P<minute>\S+)\s+(?P<hour>\S+)\s+(?P<dom>\S+)\s+(?P<month>\S+)\s+(?P<dow>\S+)\s*$`
	CrontabScheduleSubMatchLen  = 6
	ServiceConfigGroup          = "service"
	CronConfigGroup             = "cron"
	ServiceKeyFormat            = "service.%s"
	CronKeyFormat               = "cron.%s"
)

var (
	crontabScheduleRegex = regexp.MustCompile(CrontabScheduleRegexPattern)
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
		if cronTask.CronSchedule, err = parseYamlSchedule(cronTask.Schedule); err != nil {
			wingmate.Log().Error().Msgf("parsing cron schedule: %+v", err)
			continue
		}
		crones = append(crones, cronTask)
	}

	return services, crones, nil
}

func parseYamlSchedule(input string) (schedule CronSchedule, err error) {
	var (
		parts  []string
		pSched *CronSchedule
	)

	parts = crontabScheduleRegex.FindStringSubmatch(input)
	if len(parts) != CrontabScheduleSubMatchLen {
		return schedule, fmt.Errorf("invalid schedule: %s", input)
	}

	pSched = &schedule
	if err = pSched.setField(minute, parts[1]); err != nil {
		return schedule, fmt.Errorf("error parsing Minute field: %w", err)
	}

	if err = pSched.setField(hour, parts[2]); err != nil {
		return schedule, fmt.Errorf("error parsing Hour field: %w", err)
	}

	if err = pSched.setField(dom, parts[3]); err != nil {
		return schedule, fmt.Errorf("error parsing Day of Month field: %w", err)
	}

	if err = pSched.setField(month, parts[4]); err != nil {
		return schedule, fmt.Errorf("error parsing Month field: %w", err)
	}

	if err = pSched.setField(dow, parts[5]); err != nil {
		return schedule, fmt.Errorf("error parsing Day of Week field: %w", err)
	}

	return
}

func (c *CronSchedule) setField(field cronField, input string) error {
	var (
		fr       *fieldRange
		cField   *CronTimeSpec
		err      error
		parsed64 uint64
		parsed   uint8
		multi    []uint8
		current  uint8
		multiStr []string
	)
	switch field {
	case minute:
		fr = newRange(0, 59)
		cField = &c.Minute
	case hour:
		fr = newRange(0, 23)
		cField = &c.Hour
	case dom:
		fr = newRange(1, 31)
		cField = &c.DoM
	case month:
		fr = newRange(1, 12)
		cField = &c.Month
	case dow:
		fr = newRange(0, 6)
		cField = &c.DoW
	default:
		return errors.New("invalid cron field descriptor")
	}

	if input == "*" {
		*cField = &SpecAny{}
	} else if strings.HasPrefix(input, "*/") {
		if parsed64, err = strconv.ParseUint(input[2:], 10, 8); err != nil {
			return fmt.Errorf("error parse field %+v with input %s: %w", field, input, err)
		}

		parsed = uint8(parsed64)
		if !fr.valid(parsed) {
			return fmt.Errorf("error parse field %+v with input %s parsed to %d: invalid value", field, input, parsed)
		}
		multi = make([]uint8, 0)
		current = parsed
		for fr.valid(current) {
			multi = append(multi, current)
			current += parsed
		}

		*cField = &SpecMultiOccurrence{
			values: multi,
		}
	} else {
		multiStr = strings.Split(input, ",")
		if len(multiStr) > 1 {
			multi = make([]uint8, 0)
			for _, s := range multiStr {
				if parsed64, err = strconv.ParseUint(s, 10, 8); err != nil {
					return fmt.Errorf("error parse field %+v with input %s: %w", field, input, err)
				}

				parsed = uint8(parsed64)
				if !fr.valid(parsed) {
					return fmt.Errorf("error parse field %+v with input %s: invalid value", field, input)
				}

				multi = append(multi, parsed)
			}

			*cField = &SpecMultiOccurrence{
				values: multi,
			}
		} else {
			if parsed64, err = strconv.ParseUint(input, 10, 8); err != nil {
				return fmt.Errorf("error parse field %+v with input %s: %w", field, input, err)
			}

			parsed = uint8(parsed64)
			if !fr.valid(parsed) {
				return fmt.Errorf("error parse field %+v with input %s: invalid value", field, input)
			}

			*cField = &SpecExact{
				value: parsed,
			}
		}
	}

	return nil
}
