package renderer

import (
	"context"
	"testing"

	pb_fetcher "renderer-go/pb/fetcher"
	utils "renderer-go/utils"

	"github.com/stretchr/testify/assert"
)

type RenderTestCase struct {
	in       string
	out      string
	err      bool
	fetchErr bool // fetcherがerrorを返してくるようなな入力かどうか
}

var dummyFetchText = "DUMMY"

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
<li><a href="https://taiyosli.me" rel="nofollow">piyo</a></li>
</ul>
`,
		err:      false,
		fetchErr: false,
	},
	{
		in: `- [](https://google.com)`,
		out: `<ul>
<li><a href="https://google.com" rel="nofollow">` + dummyFetchText + `</a></li>
</ul>
`,
		err:      false,
		fetchErr: false,
	},
	{
		in: `- [](https://does.not.work)`,
		out: `<ul>
<li><a href="https://does.not.work" rel="nofollow">https://does.not.work</a></li>
</ul>
`,
		err:      false,
		fetchErr: true,
	},
	{
		in: `- [](https://example1.com)
- [](https://example2.com)
- [](https://example3.com)
`,
		out: `<ul>
<li><a href="https://example1.com" rel="nofollow">DUMMY</a></li>
<li><a href="https://example2.com" rel="nofollow">DUMMY</a></li>
<li><a href="https://example3.com" rel="nofollow">DUMMY</a></li>
</ul>
`,
		err:      false,
		fetchErr: false,
	},
	{
		in: `
%%%
fuga
%%%

hoge

%%%
piyo
%%%
`,
		out: `<div class="spoiler-container"><div class="spoiler">
fuga<br>
</div></div>
<p>hoge</p>
<div class="spoiler-container"><div class="spoiler">
piyo<br>
</div></div>
`,
		err:      false,
		fetchErr: false,
	},
}

func Test_Render(t *testing.T) {
	for _, testCase := range renderTestCases {
		var testFetcerClient pb_fetcher.FetcherClient
		if !testCase.fetchErr {
			testFetcerClient = utils.CreateTestFetcherClient(func(src string) string { return dummyFetchText })
		} else {
			testFetcerClient = utils.CreateTestFetcherClientWithError(func(src string) string { return dummyFetchText })
		}

		html, err := Render(context.Background(), testFetcerClient, testCase.in)
		if !testCase.err {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
		assert.Equal(t, testCase.out, html)
	}
}
