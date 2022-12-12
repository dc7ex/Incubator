package main

import (
	"context"
	//	"net"
	//	"strconv"
	"fmt"
	"time"

	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	//  "github.com/incubator/hack"
	"github.com/incubator/logger" 

	"365ex.art/secret/app/proto/v1/token"
	_ "365ex.art/secret/config"
	client "365ex.art/secret/core/grpc/client"
	"365ex.art/secret/core/grpc/consul"
)

func main() {

	//	grpcHost := hack.Getenv("GRPC_HOST", "0.0.0.0").(string)
	//	grpcPort := hack.Getenv("GRPC_PORT", 50051).(int)

	//consul的默认端口是8500,启动后可在本地访问localhost:8500
	//HelloService就是proto文件中定义的服务
	schema, err := consul.GenerateAndRegisterConsulResolver("127.0.0.1:8500", "GrpcServer")
	if err != nil {
		logger.Error("init consul resovler error: %v", err.Error())
		return
	}

	address := fmt.Sprintf("%s:///%s", schema, "GrpcServer")
	grpcClient, err := client.NewGrpcFindClient(address)
	if err != nil {
		logger.Error("Start grpc client error: %v", err.Error())
	}
	defer grpcClient.Close()

	for {
		c := tokenproto.NewTokenGreeterClient(grpcClient.ClientConn)
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		request := new(tokenproto.VerifyTokenRequest)
		request.Aid = wrapperspb.String("10000")
		request.Uid = wrapperspb.String("9b256b3d35170c636fbe4712c09562ba97ccaa48")
		request.AccessToken = wrapperspb.String("6d517208c992f00fe476d85d6449b04a8e07d90d3a4f2e86351101abb7071930")

		response, err := c.VerifyToken(ctx, request)
		logger.Error("response=%v, err=%v", response, err)
		logger.Error("ErrCode=%v", response.GetErrCode().GetValue())
		logger.Error("ErrMessage=%v", response.GetErrMessage().GetValue())
		logger.Error("Ok=%v", response.GetOk().GetValue())

		time.Sleep(time.Second * 2)
	}
}
