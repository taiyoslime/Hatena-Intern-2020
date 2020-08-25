package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	pb_fetcher "renderer-go/pb/fetcher"
	pb "renderer-go/pb/renderer"
)

type TestFetcherClient struct {
	Error       error
	FetcherFunc func(src string) string
}

func CreateTestFetcherlient(fetcherFunc func(src string) string) *TestFetcherClient {
	return &TestFetcherClient{
		Error:       nil,
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

type ServerRenderTestCase struct {
	in  string
	out string
}

var dummyFetchText = "DUMMY"

var serverRenderTestCases = []ServerRenderTestCase{
	{
		in: `# h1
## h2
### h3
- hoge
- fuga
- [piyo](https://taiyosli.me)
`,
		out: `<h1 id="h1">h1</h1>
<h2 id="h2">h2</h2>
<h3 id="h3">h3</h3>
<ul>
<li>hoge</li>
<li>fuga</li>
<li><a href="https://taiyosli.me">piyo</a></li>
</ul>
`,
	}, {
		in: `- [](https://google.com)`,
		out: `<ul>
<li><a href="https://google.com">` + dummyFetchText + `</a></li>
</ul>
`,
	},
}

func Test_Server_Render(t *testing.T) {
	var testFetcerClient pb_fetcher.FetcherClient = CreateTestFetcherlient(func(src string) string { return dummyFetchText })
	s := NewServer(testFetcerClient)

	for _, testCase := range serverRenderTestCases {
		src := testCase.in
		reply, err := s.Render(context.Background(), &pb.RenderRequest{Src: src})
		assert.NoError(t, err)
		assert.Equal(t, testCase.out, reply.Html)
	}
}
