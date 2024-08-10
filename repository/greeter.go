package repository

import "SomersaultCloud/bootstrap"

var dbs *bootstrap.Databases

func NewInternalApplicationConfig(e *bootstrap.Databases) {
	dbs = e
}
