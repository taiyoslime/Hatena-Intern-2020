package renderer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RenderTestCase struct {
	in  string
	out string
}

var renderTestCases = []RenderTestCase{
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
	},
}

func Test_Render(t *testing.T) {
	for _, testCase := range renderTestCases {
		html, err := Render(context.Background(), nil, testCase.in)
		assert.NoError(t, err)
		assert.Equal(t, html, testCase.out)
	}
}
