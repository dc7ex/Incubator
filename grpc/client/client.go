package client

import (
	"context"
	"net"
	"strconv"
	"time"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

//	"github.com/incubator/logger"
)

var GrpcServerConn *GrpcClient

type GrpcClient struct {
	*grpc.ClientConn
	host    string
	port    int
	address string
}

func NewGrpcClient(host string, port int) (grpcClient *GrpcClient, e error) {
	e = nil

	grpcClient = new(GrpcClient)
	grpcClient.host = host
	grpcClient.port = port

	// Set up a connection to the server.
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var dopts = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(host, strconv.Itoa(port)), grpc.WithBlock(), grpc.WithInsecure(), grpc.WithKeepaliveParams(dopts))
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return
	}

	grpcClient.ClientConn = conn

	GrpcServerConn = grpcClient

	log.Printf("Grpc client connection is successful \n")

	return
}

func NewGrpcFindClient(address string) (grpcClient *GrpcClient, e error) {
	e = nil

	grpcClient = new(GrpcClient)
	grpcClient.address = address

	// Set up a connection to the server.
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var dopts = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	conn, err := grpc.DialContext(ctx, address, grpc.WithBlock(), grpc.WithInsecure(), grpc.WithKeepaliveParams(dopts), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		return
	}

	grpcClient.ClientConn = conn

	GrpcServerConn = grpcClient

	log.Printf("Grpc client connection is successful \n")

	return
}

func (g *GrpcClient) Close() error {
	return g.Close()
}
