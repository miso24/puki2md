package inline

import (
	"strings"

	"github.com/miso24/puki2md/model"
	"github.com/miso24/puki2md/renderer"
)

type WrapRule struct {
	prefix      string
	suffix      string
	allowInline bool
	build       func([]model.InlineNode) model.InlineNode
}

var wrapRules = []WrapRule{
	WrapRule{
		prefix:      "'''",
		suffix:      "'''",
		allowInline: true,
		build: func(children []model.InlineNode) model.InlineNode {
			return model.Italic{Inlines: children}
		},
	},
	WrapRule{
		prefix:      "''",
		suffix:      "''",
		allowInline: true,
		build: func(children []model.InlineNode) model.InlineNode {
			return model.Strong{Inlines: children}
		},
	},
	WrapRule{
		prefix:      "((",
		suffix:      "))",
		allowInline: true,
		build: func(children []model.InlineNode) model.InlineNode {
			return model.FootNote{Inlines: children}
		},
	},
	WrapRule{
		prefix:      "[[",
		suffix:      "]]",
		allowInline: false,
		build: func(children []model.InlineNode) model.InlineNode {
			content := renderer.RenderInline(children[0])
			pageName := content
			target := ""
			for _, del := range []string{">", ":"} {
				if idx := strings.Index(content, del); idx != -1 {
					elems := strings.SplitN(content, del, 2)
					pageName = elems[0]
					target = elems[1]
				}
			}
			return model.PageLink{Name: pageName, Target: target}
		},
	},
	WrapRule{
		prefix:      "&ref(",
		suffix:      ");",
		allowInline: false,
		build: func(children []model.InlineNode) model.InlineNode {
			fileName := renderer.RenderInline(children[0])
			return model.Attach{FileName: fileName}
		},
	},
	WrapRule{
		prefix:      "&code_x{",
		suffix:      "};",
		allowInline: false,
		build: func(children []model.InlineNode) model.InlineNode {
			return model.RawCode{Code: renderer.RenderInline(children)}
		},
	},
}

func hasSpecialPrefix(line string) bool {
	for _, rule := range wrapRules {
		if strings.HasPrefix(line, rule.prefix) {
			return true
		}
	}
	return false
}

func Parse(line string) []model.InlineNode {
	nodes := []model.InlineNode{}
	i := 0
	for i < len(line) {
		slice := line[i:]
		match := false
		for _, rule := range wrapRules {
			if !strings.HasPrefix(slice, rule.prefix) {
				continue
			}

			end := strings.Index(slice[len(rule.prefix):], rule.suffix)
			if end == -1 {
				continue
			}

			content := slice[len(rule.prefix) : end+len(rule.suffix)]
			if rule.allowInline {
				children := Parse(content)
				nodes = append(nodes, rule.build(children))
			} else {
				children := []model.InlineNode{model.PlainText{Content: content}}
				nodes = append(nodes, rule.build(children))
			}
			i += end + len(rule.suffix) + len(rule.prefix)
			match = true
			break
		}

		if !match {
			var buf strings.Builder
			for i < len(line) && !hasSpecialPrefix(line[i:]) {
				buf.WriteByte(line[i])
				i++
			}
			nodes = append(nodes, model.PlainText{Content: buf.String()})
		}
	}
	return nodes
}
