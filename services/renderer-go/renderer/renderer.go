package renderer

import (
	"bytes"
	"context"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	pb_fetcher "renderer-go/pb/fetcher"
)

// Render は受け取った文書を HTML に変換する
func Render(ctx context.Context, fetcherClient pb_fetcher.FetcherClient, src string) (string, error) {

	reply, err := fetcherClient.Fetch(ctx, &pb_fetcher.FetchRequest{Url: "https://google.com"})
	if err != nil {
		return "error", err
	}
	
	parser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
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
	return html + reply.Title, nil
}
