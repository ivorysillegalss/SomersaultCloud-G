package bootstrap

import (
	"SomersaultCloud/constant/sys"
	"SomersaultCloud/domain"
)

func NewChannel() *Channels {
	//TODO channel类型的增多可以在channel结构体中增加，并在此处初始化
	return &Channels{RpcRes: make(chan *domain.GenerationResponse, sys.GenerationResponseChannelBuffer),
		StreamRpcRes: make(chan *domain.GenerationResponse, sys.StreamGenerationResponseChannelBuffer),
		Stop:         make(chan bool)}
}
