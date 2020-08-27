package renderer

import (
	"bytes"
	"context"
	pb_fetcher "renderer-go/pb/fetcher"
	"sync"

	lextension "renderer-go/renderer/extension"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type lookUpBlankTitle struct {
	dest []string
}

func (l *lookUpBlankTitle) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if node, ok := node.(*ast.Link); ok && entering && node.ChildCount() == 0 {
			l.dest = append(l.dest, string(node.Destination))
		}
		return ast.WalkContinue, nil
	})
}

type autoTitleLinker struct {
	titleMap *map[string]string
}

func (a *autoTitleLinker) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if node, ok := node.(*ast.Link); ok && entering && node.ChildCount() == 0 {
			node.AppendChild(node, ast.NewString([]byte((*(a.titleMap))[string(node.Destination)])))
		}
		return ast.WalkContinue, nil
	})
}

type fetcherStruct struct {
	ctx           context.Context
	fetcherClient pb_fetcher.FetcherClient
}

func createFetcher(ctx context.Context, fetcherClient pb_fetcher.FetcherClient) *fetcherStruct {
	return &fetcherStruct{
		ctx:           ctx,
		fetcherClient: fetcherClient,
	}
}

func (f *fetcherStruct) fetch(url string) string {
	reply, err := f.fetcherClient.Fetch(f.ctx, &pb_fetcher.FetchRequest{Url: url})
	if err != nil {
		return url // titleが得られなかった場合はurlをそのまま返す
	}
	return reply.Title
}

// Render は受け取った文書を HTML に変換する
func Render(ctx context.Context, fetcherClient pb_fetcher.FetcherClient, src string) (string, error) {

	var lookup = &lookUpBlankTitle{}

	preParser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithASTTransformers(
				util.Prioritized(lookup, 99),
			),
		),
	)
	var prebuf bytes.Buffer
	if err := preParser.Convert([]byte(src), &prebuf); err != nil {
		return src, err
	}

	fetcher := createFetcher(ctx, fetcherClient)

	titleMap := make(map[string]string)

	var wg sync.WaitGroup
	for _, url := range lookup.dest {
		wg.Add(1)
		go func(url string) {
			titleMap[url] = fetcher.fetch(url)
			wg.Done()
		}(url)
	}
	wg.Wait()

	linker := &autoTitleLinker{&titleMap}

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

	html := buf.String()
	return html, nil
}
