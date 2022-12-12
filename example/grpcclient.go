package main

import (
	"context"
	"time"

	"github.com/incubator/hack"
	"github.com/incubator/logger"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"

	_ "365ex.art/secret/config"
	client "365ex.art/secret/core/grpc/client"
	"365ex.art/secret/app/proto/v1/token"
)

func main() {

	grpcHost := hack.Getenv("GRPC_HOST", "0.0.0.0").(string)
	grpcPort := hack.Getenv("GRPC_PORT", 50051).(int)

	grpcClient, err := client.NewGrpcClient(grpcHost, grpcPort)
	if err != nil {
		logger.Error("Start grpc client error: %v", err.Error())
	}
	defer grpcClient.Close()

	c := tokenproto.NewTokenGreeterClient(grpcClient.ClientConn)
	for {
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
