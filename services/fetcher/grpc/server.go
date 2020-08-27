package grpc

import (
	"context"

	fetcher "fetcher/fetch"
	pb "fetcher/pb/fetcher"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Server は pb.RendererServer に対する実装
type Server struct {
	pb.UnimplementedFetcherServer
	healthpb.UnimplementedHealthServer
	cacheClient fetcher.CacheClient
}

// NewServer は gRPC サーバーを作成する
func NewServer(cacheClient fetcher.CacheClient) *Server {
	return &Server{
		cacheClient: cacheClient,
	}
}

func (s *Server) Fetch(ctx context.Context, in *pb.FetchRequest) (*pb.FetchReply, error) {
	title, err := fetcher.Fetch(ctx, s.cacheClient, in.Url)
	if err != nil {
		return nil, err
	}
	return &pb.FetchReply{Title: title}, nil
}
