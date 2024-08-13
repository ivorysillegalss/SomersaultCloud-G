package tokenutil

import "SomersaultCloud/bootstrap"

var env *bootstrap.Env

func NewInternalApplicationConfig(e *bootstrap.Env) {
	env = e
}
