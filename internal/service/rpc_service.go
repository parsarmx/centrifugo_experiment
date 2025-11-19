package service

// dto "golang_template/api/dtos"

type RpcServiceService interface {
}

type rpcServiceService struct{}

func NewRpcServiceService() RpcServiceService { // repo repositories.ExampleRepository
	return &rpcServiceService{}
}
