package pool

import (
	"github.com/panjf2000/ants/v2"
	"sync"
)

type Factory interface {
}

type Config struct {
	PoolType int
	PoolName string
	Task     func(i interface{})
	Pool     *ants.PoolWithFunc
	//确保只执行一次
	Once sync.Once
}

type Pool struct {
	*ants.PoolWithFunc
	*ants.Pool
}
