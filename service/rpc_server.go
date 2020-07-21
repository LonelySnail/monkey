package service

type IRpcServer interface {
}

type RpcServer struct {
}

type IRpcClient interface {
	Call() error
	CallNR() error
}

type RpcClient struct {
}
