package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gitea.suyono.dev/suyono/wingmate"
)

type CronExactSpec interface {
	CronTimeSpec
	Value() uint8
}

type CronMultipleOccurrenceSpec interface {
	CronTimeSpec
	Values() []uint8
}

type CronTimeSpec interface {
	Type() wingmate.CronTimeType
	Match(uint8) bool
}

type Cron struct {
	minute  CronTimeSpec
	hour    CronTimeSpec
	dom     CronTimeSpec
	month   CronTimeSpec
	dow     CronTimeSpec
	command string
	lastRun time.Time
	hasRun  bool
}

type cronField int

const (
	CrontabEntryRegex  = `^\s*(?P<minute>\S+)\s+(?P<hour>\S+)\s+(?P<dom>\S+)\s+(?P<month>\S+)\s+(?P<dow>\S+)\s+(?P<command>\S.*\S)\s*$`
	CrontabSubmatchLen = 7

	minute cronField = iota
	hour
	dom
	month
	dow
)

func readCrontab(path string) ([]*Cron, error) {
	var (
		file    *os.File
		err     error
		scanner *bufio.Scanner
		line    string
		re      *regexp.Regexp
		parts   []string
		retval  []*Cron
	)

	if re, err = regexp.Compile(CrontabEntryRegex); err != nil {
		return nil, err
	}

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

		parts = re.FindStringSubmatch(line)
		if len(parts) != CrontabSubmatchLen {
			wingmate.Log().Error().Msgf("invalid entry %s", line)
			continue
		}

		c := &Cron{
			hasRun: false,
		}
		if err = c.setField(minute, parts[1]); err != nil {
			wingmate.Log().Error().Msgf("error parsing minute field %#v", err)
			continue
		}

		if err = c.setField(hour, parts[2]); err != nil {
			wingmate.Log().Error().Msgf("error parsing hour field %#v", err)
			continue
		}

		if err = c.setField(dom, parts[3]); err != nil {
			wingmate.Log().Error().Msgf("error parsing day of month field %#v", err)
			continue
		}

		if err = c.setField(month, parts[4]); err != nil {
			wingmate.Log().Error().Msgf("error parsing month field %#v", err)
			continue
		}

		if err = c.setField(dow, parts[5]); err != nil {
			wingmate.Log().Error().Msgf("error parsing day of week field %#v", err)
			continue
		}

		c.command = parts[6]

		retval = append(retval, c)
	}

	return retval, nil
}

func (c *Cron) Command() string {
	return c.command
}

func (c *Cron) TimeToRun(now time.Time) bool {
	if !c.hasRun {
		c.lastRun = now
		c.hasRun = true
		return true
	}

	if now.Sub(c.lastRun) <= time.Minute && now.Minute() == c.lastRun.Minute() {
		return false
	}

	if c.minute.Match(uint8(now.Minute())) &&
		c.hour.Match(uint8(now.Hour())) &&
		c.dom.Match(uint8(now.Day())) &&
		c.month.Match(uint8(now.Month())) &&
		c.dow.Match(uint8(now.Weekday())) {
		c.lastRun = now
		return true
	}

	return false
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
		cField = &c.minute
	case hour:
		fr = newRange(0, 23)
		cField = &c.hour
	case dom:
		fr = newRange(1, 31)
		cField = &c.dom
	case month:
		fr = newRange(1, 12)
		cField = &c.month
	case dow:
		fr = newRange(0, 6)
		cField = &c.dow
	default:
		return errors.New("invalid cron field descriptor")
	}

	if input == "*" {
		*cField = &specAny{}
	} else if strings.HasPrefix(input, "*/") {
		if parsed64, err = strconv.ParseUint(input[2:], 10, 8); err != nil {
			return fmt.Errorf("error parse field %#v with input %s: %w", field, input, err)
		}

		parsed = uint8(parsed64)
		if fr.valid(parsed) {
			return fmt.Errorf("error parse field %#v with input %s: invalid value", field, input)
		}
		multi = make([]uint8, 0)
		current = parsed
		for fr.valid(current) {
			multi = append(multi, current)
			current += parsed
		}

		*cField = &specMultiOccurrence{
			values: multi,
		}
	} else {
		multiStr = strings.Split(input, ",")
		if len(multiStr) > 1 {
			multi = make([]uint8, 0)
			for _, s := range multiStr {
				if parsed64, err = strconv.ParseUint(s, 10, 8); err != nil {
					return fmt.Errorf("error parse field %#v with input %s: %w", field, input, err)
				}

				parsed = uint8(parsed64)
				if fr.valid(parsed) {
					return fmt.Errorf("error parse field %#v with input %s: invalid value", field, input)
				}

				multi = append(multi, parsed)
			}

			*cField = &specMultiOccurrence{
				values: multi,
			}
		} else {
			if parsed64, err = strconv.ParseUint(input, 10, 8); err != nil {
				return fmt.Errorf("error parse field %#v with input %s: %w", field, input, err)
			}

			parsed = uint8(parsed64)
			if fr.valid(parsed) {
				return fmt.Errorf("error parse field %#v with input %s: invalid value", field, input)
			}

			*cField = &specExact{
				value: parsed,
			}
		}
	}

	return nil
}

type specAny struct{}

func (a *specAny) Type() wingmate.CronTimeType {
	return wingmate.Any
}

func (a *specAny) Match(u uint8) bool {
	return true
}

type specExact struct {
	value uint8
}

func (e *specExact) Type() wingmate.CronTimeType {
	return wingmate.Exact
}

func (e *specExact) Match(u uint8) bool {
	return u == e.value
}

func (e *specExact) Value() uint8 {
	return e.value
}

type specMultiOccurrence struct {
	values []uint8
}

func (m *specMultiOccurrence) Type() wingmate.CronTimeType {
	return wingmate.MultipleOccurrence
}

func (m *specMultiOccurrence) Match(u uint8) bool {
	for _, v := range m.values {
		if v == u {
			return true
		}
	}

	return false
}

func (m *specMultiOccurrence) Values() []uint8 {
	out := make([]uint8, len(m.values))
	copy(out, m.values)
	return out
}
