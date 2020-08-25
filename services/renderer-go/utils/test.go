package utils 

import (
	"context"
	"google.golang.org/grpc"
	"errors"
	pb_fetcher "renderer-go/pb/fetcher"
)


type TestFetcherClient struct {
	Error       error
	FetcherFunc func(src string) string
}

func CreateTestFetcherClient(fetcherFunc func(src string) string) *TestFetcherClient {
	return &TestFetcherClient{
		Error:       nil,
		FetcherFunc: fetcherFunc,
	}
}

func CreateTestFetcherClientWithError(fetcherFunc func(src string) string) *TestFetcherClient {
	return &TestFetcherClient{
		Error:       errors.New(""),
		FetcherFunc: fetcherFunc,
	}
}

func (c *TestFetcherClient) Fetch(ctx context.Context, in *pb_fetcher.FetchRequest, opts ...grpc.CallOption) (*pb_fetcher.FetchReply, error) {
	if c.Error != nil {
		return nil, c.Error
	}
	title := in.Url
	if c.FetcherFunc != nil {
		title = c.FetcherFunc(in.Url)
	}
	return &pb_fetcher.FetchReply{Title: title}, nil
}
