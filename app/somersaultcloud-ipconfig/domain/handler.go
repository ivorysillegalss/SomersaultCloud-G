package domain

type Dispatcher interface {
	Handle()
	Do(ctx *IpConfContext) []*EndPort
}

type DataHandler interface {
	Handle()
}
