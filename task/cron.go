package task

import (
	wminit "gitea.suyono.dev/suyono/wingmate/init"
	"time"
)

type CronSchedule struct {
	Minute CronTimeSpec
	Hour   CronTimeSpec
	DoM    CronTimeSpec
	Month  CronTimeSpec
	DoW    CronTimeSpec
}

type CronTimeSpec interface {
	Match(uint8) bool
}

type CronAnySpec struct {
}

func NewCronAnySpec() *CronAnySpec {
	return &CronAnySpec{}
}

func (cas *CronAnySpec) Match(u uint8) bool {
	return true
}

type CronExactSpec struct {
	value uint8
}

func NewCronExactSpec(v uint8) *CronExactSpec {
	return &CronExactSpec{
		value: v,
	}
}

func (ces *CronExactSpec) Match(u uint8) bool {
	return u == ces.value
}

type CronMultiOccurrenceSpec struct {
	values []uint8
}

func NewCronMultiOccurrenceSpec(v ...uint8) *CronMultiOccurrenceSpec {
	retval := &CronMultiOccurrenceSpec{}
	if len(v) > 0 {
		retval.values = make([]uint8, len(v))
		copy(retval.values, v)
	}

	return retval
}

func (cms *CronMultiOccurrenceSpec) Match(u uint8) bool {
	for _, v := range cms.values {
		if v == u {
			return true
		}
	}

	return false
}

type Cron struct {
	CronSchedule
	name    string
	command []string
	lastRun time.Time
	hasRun  bool //NOTE: make sure initialised as false
}

func (c *Cron) Name() string {
	return c.name
}

func (c *Cron) Command() []string {
	retval := make([]string, len(c.command))
	copy(retval, c.command)
	return retval
}

func (c *Cron) Environ() []string {
	panic("not implemented")
	return nil
}

func (c *Cron) Setsid() bool {
	panic("not implemented")
	return false
}

func (c *Cron) UserGroup() wminit.UserGroup {
	panic("not implemented")
	return nil
}

func (c *Cron) Background() bool {
	panic("not implemented")
	return false
}

func (c *Cron) WorkingDir() string {
	panic("not implemented")
	return ""
}

func (c *Cron) Status() wminit.TaskStatus {
	panic("not implemented")
	return nil
}

func (c *Cron) TimeToRun(now time.Time) bool {
	if c.Minute.Match(uint8(now.Minute())) &&
		c.Hour.Match(uint8(now.Hour())) &&
		c.DoM.Match(uint8(now.Day())) &&
		c.Month.Match(uint8(now.Month())) &&
		c.DoW.Match(uint8(now.Weekday())) {

		if c.hasRun {
			if now.Sub(c.lastRun) <= time.Minute && now.Minute() == c.lastRun.Minute() {
				return false
			} else {
				c.lastRun = now
				return true
			}
		} else {
			c.lastRun = now
			c.hasRun = true
			return true
		}
	}

	return false
}
