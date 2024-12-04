package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-chat/constant/sys"
	"SomersaultCloud/app/somersaultcloud-chat/infrastructure/pool"
	"github.com/panjf2000/ants/v2"
)

// TODO 当线程池的类型多起来之后
//
//	维护本地哈希缓存 保存poolType常量值和对应池子间的映射
func NewPoolFactory() *PoolsFactory {
	//TODO 将线程池参数信息可配置化
	p := make(map[int]*pool.Pool, sys.GoRoutinePoolTypesAmount)
	//TODO 写入线程池相关信息 & 初始化线程池 记得defer poo.release!!

	// 初始化ExecuteRpcGoRoutinePool的线程池

	defaultPool, _ := ants.NewPool(sys.DefaultPoolGoRoutineAmount) // 1000 可以是您期望的 goroutine 数量
	//ants.NewPoolWithFunc(sys.DefaultPoolGoRoutineAmount)
	p[sys.ExecuteRpcGoRoutinePool] = &pool.Pool{Pool: defaultPool}

	return &PoolsFactory{Pools: p}
}
