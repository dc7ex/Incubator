package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/incubator/logger"
)

// StreamGSError500 捕捉流式代码致命错误
func middlewareStreamGSError500() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {

		defer func() {
			if err := recover(); err != nil {
				ctx := stream.Context()
				p, _ := peer.FromContext(ctx)
				//打印错误堆栈信息
				logger.Error("StreamGSError500 ip=%v, Err: %v", p.Addr.String(), err)
			}
		}()

		err = handler(srv, stream)
		return err
	}
}

// UnaryGSError500 捕捉简单代码致命错误
func middlewareUnaryGSError500() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {

		defer func() {
			if err := recover(); err != nil {
				p, _ := peer.FromContext(ctx)
				//打印错误堆栈信息
				logger.Error("UnaryGSError500 ip=%v, Err: %v", p.Addr.String(), err)
			}
		}()

		resp, err := handler(ctx, req)
		return resp, err
	}
}
