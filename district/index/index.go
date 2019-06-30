package index

import (
	"context"
	"net"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	pb "github.com/leeif/mercury/district/index/index"
	"google.golang.org/grpc"
)

type IndexNode struct {
	logger log.Logger
}

func (i IndexNode) RegisterMember(ctx context.Context, in *pb.RegisterMemberRequest) (*pb.RegisterMemberReply, error) {
	return &pb.RegisterMemberReply{Message: "Hello " + in.Ip}, nil
}

func (i IndexNode) Start() {
	address := ":8080"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		level.Error(i.logger).Log("error", err)
	}
	s := grpc.NewServer()
	pb.RegisterIndexNodeServer(s, i)
	level.Info(i.logger).Log("msg", "index server is listening at "+address)
	if err := s.Serve(lis); err != nil {
		level.Error(i.logger).Log("error", err)
	}
}

func NewIndexNode(logger log.Logger) IndexNode {
	log.With(logger)
	indexNode := IndexNode{
		logger: logger,
	}
	return indexNode
}
