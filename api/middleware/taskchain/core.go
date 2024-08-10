package taskchain

type TaskContext struct {
	businessType    string
	businessCode    int
	taskContextData TaskContextData
	//T any 特定类型的数据
	Exception bool
}

// TaskContextData 此处定义责任链节点所使用的共同数据
type TaskContextData struct {
}

type TaskContextResponse struct {
	Message string
	Code    int
}

type TaskNodeModel interface {
	Execute(ctx *TaskContext)
}

// ExecuteChain 执行责任链
func (t *TaskContext) ExecuteChain(handlers ...TaskNodeModel) {
	for _, handler := range handlers {

		if t.Exception {
			// 具体的错误原因类型在 实现类中包装
			return
		}
		handler.Execute(t)
	}
}

// TODO 目前是即用即装配链子 可在此基础上二次封装template
// 			kv形式存储		 k for 业务类型 & v for 链子配置 形成工厂类策略化等等乱七八糟的
// 			这个以后完善 没有实际链子太抽象了
