package main

import (
	"context"
	"strconv"
	"time"
	"net"

	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	"github.com/incubator/hack"
	"github.com/incubator/logger"

	_ "365ex.art/secret/config"
    pool "365ex.art/secret/core/grpc/client/pool"
	"365ex.art/secret/app/proto/v1/token"
)

func main() {

	grpcHost := hack.Getenv("GRPC_HOST", "0.0.0.0").(string)
	grpcPort := hack.Getenv("GRPC_PORT", 50051).(int)

	grepPool, err := pool.New(net.JoinHostPort(grpcHost, strconv.Itoa(grpcPort)), pool.DefaultOptions)
	if err != nil {
		logger.Error("failed to new pool error: %v", err.Error())
		return
	}
	defer grepPool.Close()

	conn, err := grepPool.Get()
	if err != nil {
		logger.Error("failed to get conn error: %v", err.Error())
		return
	}
	defer conn.Close()

	c := tokenproto.NewTokenGreeterClient(conn.Value())
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
