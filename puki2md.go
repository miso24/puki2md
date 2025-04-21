package puki2md

import (
	"io"
	"strings"

	"github.com/miso24/puki2md/block"
	"github.com/miso24/puki2md/renderer"
)

func Convert(r io.Reader) string {
	p := block.NewParser(r)
	blocks := p.Parse()

	var s strings.Builder
	for _, block := range blocks {
		md := renderer.RenderBlock(block)
		s.WriteString(md + "\n")
	}
	return s.String()
}
