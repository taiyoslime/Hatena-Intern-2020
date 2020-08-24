package renderer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Render(t *testing.T) {
	src := `
# h1
## h2
### h3
- hoge
- fuga
- [piyo](https://taiyosli.me)
`
	html, err := Render(context.Background(), src)
	assert.NoError(t, err)
	assert.Equal(t, `<h1 id="h1">h1</h1>
<h2 id="h2">h2</h2>
<h3 id="h3">h3</h3>
<ul>
<li>hoge</li>
<li>fuga</li>
<li><a href="https://taiyosli.me">piyo</a></li>
</ul>
`, html)
}
