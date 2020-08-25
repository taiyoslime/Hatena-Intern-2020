package grpc

import (
	"context"
	"testing"

	pb "fetcher/pb/fetcher"
	"github.com/stretchr/testify/assert"
)

func Test_Server_Fetch(t *testing.T) {
	s := NewServer()
	url := "https://hatenablog.com/"
	reply, err := s.Fetch(context.Background(), &pb.FetchRequest{Url: url})
	assert.NoError(t, err)
	assert.Equal(t, "はてなブログ", reply.Title)
}
