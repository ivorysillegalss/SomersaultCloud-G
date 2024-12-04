package bootstrap

import (
	"SomersaultCloud/app/somersaultcloud-chat/constant/sys"
	"SomersaultCloud/app/somersaultcloud-chat/domain"
)

func NewChannel() *Channels {
	//TODO channel类型的增多可以在channel结构体中增加，并在此处初始化
	return &Channels{RpcRes: make(chan *domain.GenerationResponse, sys.GenerationResponseChannelBuffer),
		StreamRpcRes: make(chan *domain.GenerationResponse, sys.StreamGenerationResponseChannelBuffer),
		Stop:         make(chan bool)}
}
