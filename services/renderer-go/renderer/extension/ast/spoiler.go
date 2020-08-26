package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

type SpoilerBlock struct {
	gast.BaseBlock
}

func CreateSpoilerBlock() *SpoilerBlock {
	return &SpoilerBlock{}
}

func (b *SpoilerBlock) Dump(source []byte, level int) {
	gast.DumpHelper(b, source, level, nil, nil)
}

var KindSpoilerBlock = gast.NewNodeKind("SpoilerBlock")

func (b *SpoilerBlock) Kind() gast.NodeKind {
	return KindSpoilerBlock
}

func CreateDefinitionList(offset int, para *gast.Paragraph) *SpoilerBlock {
	return &SpoilerBlock{}
}
