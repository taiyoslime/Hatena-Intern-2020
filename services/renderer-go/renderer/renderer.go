package renderer

import (
	"bytes"
	"context"
	pb_fetcher "renderer-go/pb/fetcher"
	"sync"

	lextension "renderer-go/renderer/extension"

	"github.com/microcosm-cc/bluemonday"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type autoTitleLinker struct {
	ctx           context.Context
	fetcherClient pb_fetcher.FetcherClient
}

func createAutoTitleLinker(ctx context.Context, fetcherClient pb_fetcher.FetcherClient) *autoTitleLinker {
	return &autoTitleLinker{
		ctx:           ctx,
		fetcherClient: fetcherClient,
	}
}

func (a *autoTitleLinker) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	var dest []string

	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if node, ok := node.(*ast.Link); ok && entering && node.ChildCount() == 0 {
			dest = append(dest, string(node.Destination))
		}
		return ast.WalkContinue, nil
	})

	var titleMap sync.Map

	var wg sync.WaitGroup
	for _, url := range dest {
		wg.Add(1)
		go func(url string) {
			if _, ok := titleMap.Load(url); !ok {
				title := a.fetch(url)
				titleMap.Store(url, title)
			}
			wg.Done()
		}(url)
	}
	wg.Wait()

	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if node, ok := node.(*ast.Link); ok && entering && node.ChildCount() == 0 {
			title, _ := titleMap.Load(string(node.Destination))
			node.AppendChild(node, ast.NewString([]byte(title.(string))))
		}
		return ast.WalkContinue, nil
	})
}

func (a *autoTitleLinker) fetch(url string) string {
	reply, err := a.fetcherClient.Fetch(a.ctx, &pb_fetcher.FetchRequest{Url: url})
	if err != nil {
		return url // titleが得られなかった場合はurlをそのまま返す
	}
	return reply.Title
}

// Render は受け取った文書を HTML に変換する
func Render(ctx context.Context, fetcherClient pb_fetcher.FetcherClient, src string) (string, error) {

	var p = bluemonday.UGCPolicy()
	p.AllowStandardAttributes()
	p.AllowStyling()

	linker := &autoTitleLinker{ctx, fetcherClient}

	parser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithASTTransformers(
				util.Prioritized(linker, 99),
			),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			lextension.SpoilerBlock(),
		),
	)

	var buf bytes.Buffer
	if err := parser.Convert([]byte(src), &buf); err != nil {
		return src, err
	}

	html := p.Sanitize(buf.String())

	return html, nil
}
