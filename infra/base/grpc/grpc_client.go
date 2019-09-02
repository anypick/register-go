package basegrpc

import (
	"google.golang.org/grpc"
	"register-go/infra"
	"time"
)

var grpcConn map[string]*grpc.ClientConn

type GrpcClientStarter struct {
	infra.BaseStarter
}

func GetGrpcConn(appName string) *grpc.ClientConn {
	return grpcConn[appName]
}

func (g *GrpcClientStarter) Start(ctx infra.StarterContext) {
	var (
		config = ctx.Yaml().GrpcConfig
		conn   *grpc.ClientConn
		err    error
	)
	grpcConn = make(map[string]*grpc.ClientConn, len(config.Clients))
	clients := config.Clients
	for _, client := range clients {
		if conn, err = grpc.Dial(client.Addr, grpc.WithInsecure()); err != nil {
			flag := make(chan struct{})
			// 出现错误重试， 每秒钟重试一次
			go func() {
				for {
					conn, err = grpc.Dial(client.Addr, grpc.WithInsecure())
					if err == nil {
						close(flag)
					}
					select {
					case <-time.After(time.Second * 1):
						conn, err = grpc.Dial(client.Addr, grpc.WithInsecure())
						if err == nil {
							close(flag)
						}
					case <-flag:
						// 跳出循环
						goto OUT
					}
				}
			OUT:
			}()
		}

		grpcConn[client.AppName] = conn
	}
}

// 有可能会被阻塞
func (g *GrpcClientStarter) StartBlocking() bool {
	return true
}
