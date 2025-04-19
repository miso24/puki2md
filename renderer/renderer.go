package renderer

import (
	"fmt"
	"strings"

	"github.com/miso24/puki2md/model"
)

func RenderInline(node model.InlineNode) string {
	var s strings.Builder
	switch n := node.(type) {
	case model.PlainText:
		s.WriteString(n.Content)
	case model.Attach:
        // TODO: implement
	case model.Italic:
		s.WriteString(fmt.Sprintf("*%s*", joinInlines(n.Inlines)))
	case model.Strong:
		s.WriteString(fmt.Sprintf("**%s**", joinInlines(n.Inlines)))
	case model.RawCode:
		s.WriteString(fmt.Sprintf("`%s`", n.Code))
	case model.FootNote:
		s.WriteString(fmt.Sprintf("[^%s]", joinInlines(n.Inlines)))
	case model.PageLink:
		if n.Target == "" {
			s.WriteString(fmt.Sprintf("[%s](%s)", n.Name, n.Name))
		} else {
			s.WriteString(fmt.Sprintf("[%s](%s)", n.Name, n.Target))
		}
	}
	return s.String()
}

func renderTableRow(row model.TableRowNode) string {
	var s strings.Builder
	s.WriteByte('|')
	for _, cell := range row.Cells {
		s.WriteString(joinInlines(cell.Inlines) + "|")
	}
	return s.String()
}

func RenderBlock(node model.BlockNode) string {
	var s strings.Builder
	switch n := node.(type) {
	case model.ParagraphNode:
		s.WriteString(joinInlines(n.Inlines))
	case model.HeadingNode:
		s.WriteString(strings.Repeat("#", n.Level) + " ")
		s.WriteString(joinInlines(n.Inlines))
	case model.ListNode:
		s.WriteString(strings.Repeat(" ", (n.Level-1)*2))
		if n.Ordered {
			s.WriteString("1. ")
		} else {
			s.WriteString("- ")
		}
		s.WriteString(joinInlines(n.Inlines))
	case model.DefinitionListNode:
		s.WriteString("<dl>\n")
		for _, item := range n.Items {
			s.WriteString(fmt.Sprintf("<dt>%s</dt>\n", joinInlines(item.Term)))
			s.WriteString(fmt.Sprintf("<dd>%s</dd>\n", joinInlines(item.Desc)))
		}
		s.WriteString("</dl>")
	case model.PreformattedNode:
		s.WriteString("```\n")
		s.WriteString(strings.Join(n.Lines, "\n") + "\n")
		s.WriteString("```")
	case model.TableNode:
		for i, row := range n.Rows {
			s.WriteString(renderTableRow(row) + "\n")
			if i == 0 {
				s.WriteString(strings.Repeat("|---", len(row.Cells)))
				s.WriteString("|\n")
			}
		}
	}
	return s.String()
}

func joinInlines(inlines []model.InlineNode) string {
	var s strings.Builder
	for _, inline := range inlines {
		s.WriteString(RenderInline(inline))
	}
	return s.String()
}
