package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	pb "renderer-go/pb/renderer"
)

func Test_Server_Render(t *testing.T) {
	s := NewServer()
	src := `
# h1
## h2
### h3
- hoge
- fuga
- [piyo](https://taiyosli.me)
`
	reply, err := s.Render(context.Background(), &pb.RenderRequest{Src: src})
	assert.NoError(t, err)
	assert.Equal(t, `<h1 id="h1">h1</h1>
<h2 id="h2">h2</h2>
<h3 id="h3">h3</h3>
<ul>
<li>hoge</li>
<li>fuga</li>
<li><a href="https://taiyosli.me">piyo</a></li>
</ul>
`, reply.Html)
}
