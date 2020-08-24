package renderer

import (
	"bytes"
	"context"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Render は受け取った文書を HTML に変換する
func Render(ctx context.Context, src string) (string, error) {

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
	return html, nil
}
