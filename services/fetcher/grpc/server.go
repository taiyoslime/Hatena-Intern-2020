package grpc

import (
	"context"

	pb "fetcher/pb/fetcher"
	"fetcher/fetch"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Server は pb.RendererServer に対する実装
type Server struct {
	pb.UnimplementedFetcherServer
	healthpb.UnimplementedHealthServer
}

// NewServer は gRPC サーバーを作成する
func NewServer() *Server {
	return &Server{}
}


func (s *Server) Fetch(ctx context.Context, in *pb.FetchRequest) (*pb.FetchReply, error) {
	title, err := fetcher.Fetch(ctx, in.Url)
	if err != nil {
		return nil, err
	}
	return &pb.FetchReply{Title: title}, nil
}
