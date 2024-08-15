package task

import "SomersaultCloud/bootstrap"

var env *bootstrap.Env
var poolFactory *bootstrap.PoolsFactory

func NewUseCaseApplicationConfig(e *bootstrap.Env, p *bootstrap.PoolsFactory) {
	env = e
	poolFactory = p
}
