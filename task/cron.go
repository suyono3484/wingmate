package task

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"gitea.suyono.dev/suyono/wingmate"

	wminit "gitea.suyono.dev/suyono/wingmate/init"
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
	cronScheduleString string
	name               string
	command            []string
	cmdLine            []string
	environ            []string
	setsid             bool
	workingDir         string
	lastRun            time.Time
	hasRun             bool //NOTE: make sure initialised as false
	config             config
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

func (c *CronTask) SetSchedule(scheduleStr string, schedule CronSchedule) *CronTask {
	c.cronScheduleString = scheduleStr
	c.CronSchedule = schedule
	return c
}

func (c *CronTask) SetConfig(config config) *CronTask {
	c.config = config
	return c
}

func (c *CronTask) Equals(another *CronTask) bool {
	if another == nil {
		return false
	}

	type toCompare struct {
		Name       string
		Command    string
		Arguments  []string
		Environ    []string
		Setsid     bool
		UserGroup  string
		WorkingDir string
		Schedule   string
	}

	cmpStruct := func(p *CronTask) ([]byte, error) {
		s := &toCompare{
			Name:       p.Name(),
			Command:    p.Command(),
			Arguments:  p.Arguments(),
			Environ:    p.Environ(),
			Setsid:     p.Setsid(),
			UserGroup:  p.UserGroup().String(),
			WorkingDir: p.WorkingDir(),
			Schedule:   p.cronScheduleString,
		}

		return json.Marshal(s)
	}

	var (
		err                error
		ours, theirs       []byte
		ourHash, theirHash [sha256.Size]byte
	)

	if ours, err = cmpStruct(c); err != nil {
		wingmate.Log().Error().Msgf("cron task equals: %+v", err)
		return false
	}
	ourHash = sha256.Sum256(ours)

	if theirs, err = cmpStruct(another); err != nil {
		wingmate.Log().Error().Msgf("cron task equals: %+v", err)
		return false
	}
	theirHash = sha256.Sum256(theirs)

	for i := 0; i < sha256.Size; i++ {
		if ourHash[i] != theirHash[i] {
			return false
		}
	}

	return true
}

func (c *CronTask) Name() string {
	return c.name
}

func (c *CronTask) UtilDepCheck() error {
	c.cmdLine = make([]string, 0)
	if c.setsid || c.UserGroup().IsSet() {
		if err := c.config.WMExecCheckVersion(); err != nil {
			return fmt.Errorf("utility dependency check: %w", err)
		}

		c.cmdLine = append(c.cmdLine, c.config.WMExecPath())

		if c.setsid {
			c.cmdLine = append(c.cmdLine, "--setsid")
		}

		if c.UserGroup().IsSet() {
			c.cmdLine = append(c.cmdLine, "--user", c.UserGroup().String())
		}

		c.cmdLine = append(c.cmdLine, "--")
	}

	c.cmdLine = append(c.cmdLine, c.command...)

	return nil
}

func (c *CronTask) Command() string {
	return c.cmdLine[0]
}

func (c *CronTask) Arguments() []string {
	if len(c.cmdLine) == 1 {
		return nil
	}

	retval := make([]string, len(c.cmdLine)-1)
	copy(retval, c.cmdLine[1:])

	return retval
}

func (c *CronTask) EnvLen() int {
	return len(c.environ)
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
