package executor

import (
	"SomersaultCloud/domain"
)

type CronExecutor struct {
	GenerationCron domain.GenerationCron
}

func (d *CronExecutor) SetupCron() {
	go d.GenerationCron.AsyncPollerGeneration()
}

func NewCronExecutor(g domain.GenerationCron) *CronExecutor {
	return &CronExecutor{GenerationCron: g}
}
