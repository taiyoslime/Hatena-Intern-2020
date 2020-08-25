package renderer

import (
	"bytes"
	"context"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	pb_fetcher "renderer-go/pb/fetcher"
)

type autoTitleLinker struct {
	ctx           context.Context
	fetcherClient pb_fetcher.FetcherClient
}

func (l *autoTitleLinker) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if node, ok := node.(*ast.Link); ok && entering && node.ChildCount() == 0 {
			node.AppendChild(node, ast.NewString([]byte(fetch(l, string(node.Destination)))))
		}
		return ast.WalkContinue, nil
	})
}

func fetch(l *autoTitleLinker, url string) string {
	reply, err := l.fetcherClient.Fetch(l.ctx, &pb_fetcher.FetchRequest{Url: url})
	if err != nil {
		return ""
	}
	return reply.Title
}

// Render は受け取った文書を HTML に変換する
func Render(ctx context.Context, fetcherClient pb_fetcher.FetcherClient, src string) (string, error) {

	parser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithASTTransformers(
				util.Prioritized(&autoTitleLinker{ctx, fetcherClient}, 99),
			),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	var buf bytes.Buffer
	if err := parser.Convert([]byte(src), &buf); err != nil {
		return src, err
	}

	html := buf.String()
	return html, nil
}
