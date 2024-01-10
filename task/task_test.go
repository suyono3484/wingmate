package task

import (
	wminit "gitea.suyono.dev/suyono/wingmate/init"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServicesV0(t *testing.T) {
	service := "/path/to/executable"
	tasks := NewTasks()
	tasks.AddV0Service(service)

	assert.Equal(t, tasks.Services()[0].Name(), service)
	assert.ElementsMatch(t, tasks.Services()[0].Command(), []string{service})
}

func TestCronV0(t *testing.T) {
	cron := "/path/to/executable"
	tasks := NewTasks()
	tasks.AddV0Cron(CronSchedule{
		Minute: NewCronAnySpec(),
		Hour:   NewCronAnySpec(),
		DoM:    NewCronAnySpec(),
		Month:  NewCronAnySpec(),
		DoW:    NewCronAnySpec(),
	}, cron)

	assert.Equal(t, tasks.Crones()[0].Name(), cron)
	assert.ElementsMatch(t, tasks.Crones()[0].Command(), []string{cron})
}

func TestTasks_List(t *testing.T) {
	tasks := NewTasks()
	tasks.services = []wminit.ServiceTask{
		&ServiceTask{
			name:    "one",
			command: []string{"/path/to/executable"},
		},
		&ServiceTask{
			name:    "two",
			command: []string{"/path/to/executable"},
		},
	}
	tasks.crones = []wminit.CronTask{
		&CronTask{
			CronSchedule: CronSchedule{
				Minute: NewCronAnySpec(),
				Hour:   NewCronAnySpec(),
				DoM:    NewCronAnySpec(),
				Month:  NewCronAnySpec(),
				DoW:    NewCronAnySpec(),
			},
			name:    "cron-one",
			command: []string{"/path/to/executable"},
		},
		&CronTask{
			CronSchedule: CronSchedule{
				Minute: NewCronAnySpec(),
				Hour:   NewCronAnySpec(),
				DoM:    NewCronAnySpec(),
				Month:  NewCronAnySpec(),
				DoW:    NewCronAnySpec(),
			},
			name:    "cron-two",
			command: []string{"/path/to/executable"},
		},
	}

	tl := tasks.List()
	tnames := make([]string, 0)
	testNames := []string{"one", "two", "cron-one", "cron-two"}

	for _, ti := range tl {
		tnames = append(tnames, ti.Name())
	}

	assert.ElementsMatch(t, testNames, tnames)
}
