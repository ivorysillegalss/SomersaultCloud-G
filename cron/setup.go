package cron

import "SomersaultCloud/bootstrap"

type DefaultExecutor struct {
	asyncService AsyncService
}

func (d *DefaultExecutor) SetupCron() {
	d.asyncService.AsyncPoller()
}

func NewExecutor() bootstrap.Executor {
	return &DefaultExecutor{}
}
