package bootstrap

import (
	"SomersaultCloud/constant/sys"
	"SomersaultCloud/infrastructure/pool"
)

// TODO 当线程池的类型多起来之后
//
//	维护本地哈希缓存 保存poolType常量值和对应池子间的映射
func NewPoolFactory() *PoolsFactory {
	//TODO 将线程池参数信息可配置化
	p := make(map[int]pool.Pool, sys.GoRoutinePoolTypesAmount)
	//TODO 写入线程池相关信息 & 初始化线程池 记得defer poo.release!!
	//p[1]=
	return &PoolsFactory{Pools: p}
}
