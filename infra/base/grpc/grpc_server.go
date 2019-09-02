package basegrpc

import (
	"google.golang.org/grpc"
	"net"
	"register-go/infra"
)

var rpc *grpc.Server

type GrpcServerStarter struct {
	infra.BaseStarter
}

func (g *GrpcServerStarter) Start(ctx infra.StarterContext) {
	config := ctx.Yaml().GrpcConfig
	lis, err := net.Listen("tcp", config.ServerConfig.Addr)
	if err != nil {
		panic(err)
	}
	rpc = grpc.NewServer()
	for _, grpcInstance := range grpcRegister.GrpcInstances {
		grpcInstance.RegisterGrpc(rpc)
	}
	if err := rpc.Serve(lis); err != nil {
		panic(err)
	}
}

func (g *GrpcServerStarter) StartBlocking() bool {
	return true
}
