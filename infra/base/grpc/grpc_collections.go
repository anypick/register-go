package basegrpc

import "google.golang.org/grpc"

/**
主要用于存储所有的rpc实例
*/
type GrpcInstance interface {
	// 将RPC实例注册到此处
	RegisterGrpc(rpc *grpc.Server)
}

type GrpcRegister struct {
	GrpcInstances []GrpcInstance
}

func (g *GrpcRegister) register(grpcInstace GrpcInstance) {
	g.GrpcInstances = append(g.GrpcInstances, grpcInstace)
}

var grpcRegister = new(GrpcRegister)

func RegisterGrpc(grpcInstance GrpcInstance) {
	grpcRegister.register(grpcInstance)
}
