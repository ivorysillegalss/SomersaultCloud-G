package cron

import "SomersaultCloud/bootstrap"

type DefaultExecutor struct {
	asyncService AsyncService
}

func (d *DefaultExecutor) SetupCron() {
	go d.asyncService.AsyncPoller()
}

func NewExecutor(service AsyncService) bootstrap.Executor {
	return &DefaultExecutor{asyncService: service}
}
