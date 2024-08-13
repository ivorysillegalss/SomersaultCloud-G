package task

import "SomersaultCloud/bootstrap"

var env *bootstrap.Env

func NewUseCaseApplicationConfig(e *bootstrap.Env) {
	env = e
}
