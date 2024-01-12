package config

import (
	"bufio"
	"errors"
	"fmt"
	"gitea.suyono.dev/suyono/wingmate"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CronTimeSpec interface {
	//Type() wingmate.CronTimeType
	//Match(uint8) bool
}

type Cron struct {
	CronSchedule
	Command string
}

type cronField int

const (
	CrontabEntryRegexPattern         = `^\s*(?P<minute>\S+)\s+(?P<hour>\S+)\s+(?P<dom>\S+)\s+(?P<month>\S+)\s+(?P<dow>\S+)\s+(?P<command>\S.*\S)\s*$`
	CrontabCommentLineRegexPattern   = `^\s*#.*$`
	CrontabCommentSuffixRegexPattern = `^\s*([^#]+)#.*$`
	CrontabSubMatchLen               = 7

	minute cronField = iota
	hour
	dom
	month
	dow
)

var (
	crontabEntryRegex         = regexp.MustCompile(CrontabEntryRegexPattern)
	crontabCommentLineRegex   = regexp.MustCompile(CrontabCommentLineRegexPattern)
	crontabCommentSuffixRegex = regexp.MustCompile(CrontabCommentSuffixRegexPattern)
)

func readCrontab(path string) ([]*Cron, error) {
	var (
		file    *os.File
		err     error
		scanner *bufio.Scanner
		line    string
		parts   []string
		retval  []*Cron
	)

	if file, err = os.Open(path); err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	retval = make([]*Cron, 0)
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()

		if crontabCommentLineRegex.MatchString(line) {
			continue
		}

		parts = crontabCommentSuffixRegex.FindStringSubmatch(line)
		if len(parts) == 2 {
			line = parts[1]
		}

		parts = crontabEntryRegex.FindStringSubmatch(line)
		if len(parts) != CrontabSubMatchLen {
			wingmate.Log().Error().Msgf("invalid entry %s", line)
			continue
		}

		c := &Cron{}
		if err = c.setField(minute, parts[1]); err != nil {
			wingmate.Log().Error().Msgf("error parsing Minute field %+v", err)
			continue
		}

		if err = c.setField(hour, parts[2]); err != nil {
			wingmate.Log().Error().Msgf("error parsing Hour field %+v", err)
			continue
		}

		if err = c.setField(dom, parts[3]); err != nil {
			wingmate.Log().Error().Msgf("error parsing Day of Month field %+v", err)
			continue
		}

		if err = c.setField(month, parts[4]); err != nil {
			wingmate.Log().Error().Msgf("error parsing Month field %+v", err)
			continue
		}

		if err = c.setField(dow, parts[5]); err != nil {
			wingmate.Log().Error().Msgf("error parsing Day of Week field %+v", err)
			continue
		}

		c.Command = parts[6]

		retval = append(retval, c)
	}

	return retval, nil
}

type fieldRange struct {
	min int
	max int
}

func newRange(min, max int) *fieldRange {
	return &fieldRange{
		min: min,
		max: max,
	}
}

func (f *fieldRange) valid(u uint8) bool {
	i := int(u)

	return i >= f.min && i <= f.max
}

func (c *Cron) setField(field cronField, input string) error {
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

type SpecAny struct{}

type SpecExact struct {
	value uint8
}

func (e *SpecExact) Value() uint8 {
	return e.value
}

type SpecMultiOccurrence struct {
	values []uint8
}

func (m *SpecMultiOccurrence) Values() []uint8 {
	out := make([]uint8, len(m.values))
	copy(out, m.values)
	return out
}
