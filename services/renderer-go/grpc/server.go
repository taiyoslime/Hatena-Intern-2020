package grpc

import (
	"context"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	pb_fetcher "renderer-go/pb/fetcher"
	pb_renderer "renderer-go/pb/renderer"
	"renderer-go/renderer"
)

// Server は pb_renderer.RendererServer に対する実装
type Server struct {
	pb_renderer.UnimplementedRendererServer
	healthpb.UnimplementedHealthServer
	fetcherClient pb_fetcher.FetcherClient
}

// NewServer は gRPC サーバーを作成する
func NewServer(fetcherClient pb_fetcher.FetcherClient) *Server {
	return &Server{
		fetcherClient: fetcherClient,
	}
}

// Render は受け取った文書を HTML に変換する
func (s *Server) Render(ctx context.Context, in *pb_renderer.RenderRequest) (*pb_renderer.RenderReply, error) {
	html, err := renderer.Render(ctx, s.fetcherClient, in.Src)
	if err != nil {
		return nil, err
	}
	return &pb_renderer.RenderReply{Html: html}, nil
}
