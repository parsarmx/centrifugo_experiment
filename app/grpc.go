package application

import (
	rpc_service "golang_template/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (a *application) InitGRPCServer(logger *zap.Logger) rpc_service.CentrifugoApiClient {
	// conn, err := grpc.NewClient("localhost:10000", grpc.WithTransportCredentials(insecure.NewCredentials())) // This is bad
	// if err != nil {
	// 	panic(err)
	// }
	// return rpc_service.NewCentrifugoApiClient(conn)
	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return rpc_service.NewCentrifugoApiClient(conn)
}
