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

type CronTask struct {
	CronSchedule
	userGroup
	name       string
	command    []string
	environ    []string
	setsid     bool
	workingDir string
	lastRun    time.Time
	hasRun     bool //NOTE: make sure initialised as false
}

func NewCronTask(name string) *CronTask {
	return &CronTask{
		name:   name,
		hasRun: false,
	}
}

func (c *CronTask) SetCommand(cmds ...string) *CronTask {
	c.command = make([]string, len(cmds))
	copy(c.command, cmds)
	return c
}

func (c *CronTask) SetEnv(envs ...string) *CronTask {
	c.environ = make([]string, len(envs))
	copy(c.environ, envs)
	return c
}

func (c *CronTask) SetFlagSetsid(flag bool) *CronTask {
	c.setsid = flag
	return c
}

func (c *CronTask) SetWorkingDir(path string) *CronTask {
	c.workingDir = path
	return c
}

func (c *CronTask) SetUser(user string) *CronTask {
	c.user = user
	return c
}

func (c *CronTask) SetGroup(group string) *CronTask {
	c.group = group
	return c
}

func (c *CronTask) SetSchedule(schedule CronSchedule) *CronTask {
	c.CronSchedule = schedule
	return c
}

func (c *CronTask) Name() string {
	return c.name
}

func (c *CronTask) Command() []string {
	retval := make([]string, len(c.command))
	copy(retval, c.command)
	return retval
}

func (c *CronTask) Environ() []string {
	retval := make([]string, len(c.environ))
	copy(retval, c.environ)
	return retval
}

func (c *CronTask) Setsid() bool {
	return c.setsid
}

func (c *CronTask) UserGroup() wminit.UserGroup {
	return &(c.userGroup)
}

func (c *CronTask) WorkingDir() string {
	return c.workingDir
}

func (c *CronTask) Status() wminit.TaskStatus {
	//TODO: implement me!
	panic("not implemented")
	return nil
}

func (c *CronTask) TimeToRun(now time.Time) bool {
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
