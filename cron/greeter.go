package cron

import (
	"SomersaultCloud/bootstrap"
)

var c *bootstrap.Channels

func NewCronApplicationConfig(channels *bootstrap.Channels) {
	c = channels
}
