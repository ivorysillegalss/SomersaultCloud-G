package handler

import (
	"SomersaultCloud/bootstrap"
)

var env *bootstrap.Env
var res *bootstrap.Channels

func NewUseCaseApplicationConfig(e *bootstrap.Env, r *bootstrap.Channels) {
	env = e
	res = r
}
