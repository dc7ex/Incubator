package server

import (
	//	"context"
	//	"encoding/hex"
	"log"
	//	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	//	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	//	"github.com/incubator/grpc/consul"
	//	"github.com/incubator/logger"
)

var GrpcServerConn *GrpcServer

type GrpcConsul struct {
	OpenServiceRegister bool
	OpenServiceFind     bool
	Name                string
	GrpcHost            string
	GrpcProt            int
	Target              string
	Interval            time.Duration
	Ttl                 int
}

type GrpcServer struct {
	*grpc.Server
	host string
	port int
	//	grpcConsul *GrpcConsul // Services Registration Find
}

func NewGrpcServer(host string, port int /*, grpcConsul *GrpcConsul*/) (grpcServer *GrpcServer, e error) {
	e = nil
	grpcServer = new(GrpcServer)
	grpcServer.host = host
	grpcServer.port = port
	//	grpcServer.grpcConsul = grpcConsul

	// grpc 调优配置
	sopts := []grpc.ServerOption{}
	sopts = append(sopts, grpc.MaxSendMsgSize(4*1024*1024))
	sopts = append(sopts, grpc.MaxRecvMsgSize(4*1024*1024))
	sopts = append(sopts, grpc.MaxConcurrentStreams(2*1024))
	sopts = append(sopts, grpc.InitialWindowSize(1*1024*1024))
	sopts = append(sopts, grpc.InitialConnWindowSize(1*1024*1024))
	sopts = append(sopts, grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     10 * time.Second, // If a client is idle for 10 seconds, send a GOAWAY
		MaxConnectionAge:      15 * time.Second, // If any connection is alive for more than 15 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}))

	sopts = append(sopts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}))

	sopts = append(sopts, grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logger.GetLogger())),
			middlewareStreamGSError500(),
			//			grpc_opentracing.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		),
	))

	sopts = append(sopts, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger.GetLogger())),
			middlewareUnaryGSError500(),
			//			grpc_opentracing.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		),
	))

	grpcServer.Server = grpc.NewServer(sopts...)

	GrpcServerConn = grpcServer

	//	log.Printf("Grpc Server running address %s prot %d\n", host, port)

	return
}

func (g *GrpcServer) Listen() (e error) {
	e = nil

	target := net.JoinHostPort(g.host, strconv.Itoa(g.port))
	listen, e := net.Listen("tcp", target)
	if e != nil {
		log.Printf("grpc server failed to listen: %v\n", e)
		return
	}

	return
}

func (g *GrpcServer) RegisterServer() (e error) {
	e = nil
	// Register reflection service on gRPC server.
	reflection.Register(g.Server)

	if e = g.Serve(listen); e != nil {
		log.Printf("failed to grpc serve: %v\n", e)
		return
	}

	return
}

func (g *GrpcServer) RegisterConsul() (e error) {
	e = nil

    /*  
        // consul 服务注册
        if g.grpcConsul != nil {
            c := g.grpcConsul
            if c.OpenServiceRegister {

                var interval time.Duration = 5 * time.Second
                var ttl int = 10
                if c.Interval > 0 {
                    interval = c.Interval
                }

                if c.Ttl > 0 {
                    ttl = c.Ttl
                }

                e = consul.Register(c.Name, c.GrpcHost, c.GrpcProt, c.Target, interval, ttl)
                if e != nil {
                    log.Printf("consul services register err: %v\n", e)
                    return
                }
            }
        }
    */

	return
}

func (g *GrpcServer) Run() (e error) {
	e = nil
	// 服务监听
	e = g.Listen()
	if e != nil {
		return
	}
	// 注册Protobuf协议
	g.RegisterProtocol()
	// Consul服务注册
	e = g.RegisterConsul()
	if e != nil {
		return
	}
	// 注册服务
	e = g.RegisterServer()

	return
}

func (g *GrpcServer) Close() {
	g.GracefulStop()
}

func (g *GrpcServer) RegisterProtocol() {}
