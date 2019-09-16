package index

import (
	"context"
	"net"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	pb "github.com/leeif/mercury/district/index/index"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type IndexNode struct {
	logger log.Logger
}

func (i IndexNode) RegisterMember(ctx context.Context, in *pb.RegisterMemberRequest) (*pb.RegisterMemberReply, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {

	}
	level.Info(i.logger).Log(pr.Addr.String())
	return &pb.RegisterMemberReply{Message: "Hello " + in.Mid}, nil
}

func (i IndexNode) GetMember(ctx context.Context, in *pb.GetMemberRequest) (*pb.GetMemberReply, error) {
	return &pb.GetMemberReply{Ip: "", Mid: in.Mid}, nil
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
	log.With(logger, "component", "index")
	indexNode := IndexNode{
		logger: logger,
	}
	return indexNode
}
