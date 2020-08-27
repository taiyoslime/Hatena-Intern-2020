package extention

import (
	"reflect"
	ast "renderer-go/renderer/extension/ast"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type spoilerParser struct {
}

func CreateSpoilerParser() parser.BlockParser {
	return &spoilerParser{}
}

func (p *spoilerParser) Trigger() []byte {
	return []byte{'%'}
}

func (p *spoilerParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if len(line) > 3 && reflect.DeepEqual(line[0:3], []byte("%%%")) {
		node := ast.CreateSpoilerBlock()
		return node, parser.NoChildren
	} else {
		return nil, parser.NoChildren
	}
}

func (p *spoilerParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if reflect.DeepEqual(line[0:3], []byte("%%%")) {
		reader.Advance(segment.Len())
		return parser.Close
	}
	node.Lines().Append(segment)
	return parser.Continue | parser.NoChildren
}

func (p *spoilerParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {

}

func (p *spoilerParser) CanInterruptParagraph() bool {
	return false
}

func (p *spoilerParser) CanAcceptIndentedLine() bool {
	return false
}

type spoilerRenderer struct {
}

func CreateSpoilerRenderer() renderer.NodeRenderer {
	return &spoilerRenderer{}
}

func (r *spoilerRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindSpoilerBlock, r.renderSpoilerBlock)
}

func (r *spoilerRenderer) renderSpoilerBlock(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {

	if entering {
		w.WriteString("<div class=\"spoiler-container\"><div class=\"spoiler\">\n")
	} else {
		w.WriteString("</div></div>\n")
	}

	return gast.WalkContinue, nil
}

type spoilerBlock struct {
}

func SpoilerBlock() *spoilerBlock {
	return &spoilerBlock{}
}

func (s *spoilerBlock) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(CreateSpoilerRenderer(), 999),
	))

	m.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(CreateSpoilerParser(), 999),
	))
}
