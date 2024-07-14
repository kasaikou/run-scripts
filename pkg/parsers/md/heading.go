package md

import (
	"unsafe"

	"github.com/yuin/goldmark/ast"
)

func bytes2Str(bytes []byte) (s string) {
	s = *(*string)(unsafe.Pointer(&bytes))
	return
}

// TextBlock is interface text block.
type TextBlock interface {
	Plaintext() string
}

// Heading contains Heading node and origin bytes data.
type Heading struct {
	node   *ast.Heading
	origin []byte
}

// NewHeading generates a Heading instance.
func NewHeading(node *ast.Heading, origin []byte) Heading {
	return Heading{node: node, origin: origin}
}

// Plaintext indicates heading text.
func (h Heading) Plaintext() string {
	return bytes2Str(h.node.Text(h.origin))
}

// Paragraph contains Paragraph node and origin bytes data.
type Paragraph struct {
	node   *ast.Paragraph
	origin []byte
}

// NewParagraph generates a Paragraph instance.
func NewParagraph(node *ast.Paragraph, origin []byte) Paragraph {
	return Paragraph{node: node, origin: origin}
}

// Plaintext indicates paragraph text.
func (p Paragraph) Plaintext() string {
	return bytes2Str(p.node.Text(p.origin))
}

// FencedCodeBlock contains FencedCodeblock node and origin bytes data.
type FencedCodeBlock struct {
	node   *ast.FencedCodeBlock
	origin []byte
}

// NewFencedCodeBlock generates a FencedCodeBlock instance.
func NewFencedCodeBlock(node *ast.FencedCodeBlock, origin []byte) FencedCodeBlock {
	return FencedCodeBlock{node: node, origin: origin}
}

// CodeBytes indicates code in the FencedCodeBlock.
func (f FencedCodeBlock) CodeBytes() []byte {
	lines := f.node.Lines()
	if lines.Len() < 1 {
		return []byte{}
	}

	return f.origin[lines.At(0).Start:lines.At(lines.Len()-1).Stop]
}

// Code indicates code in the FencedCodeBlock.
func (f FencedCodeBlock) Code() string {
	return bytes2Str(f.CodeBytes())
}

// Lang indicates language of the code in FencedCodeBlock.
func (f FencedCodeBlock) Lang() string {
	return bytes2Str(f.node.Language(f.origin))
}
