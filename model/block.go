package model

type BlockNode interface{}

type ParagraphNode struct {
	Inlines []InlineNode
}

type HeadingNode struct {
	Level   int
	Inlines []InlineNode
}

type ListNode struct {
	Level   int
	Ordered bool
	Inlines []InlineNode
}

type DefinitionListNode struct {
	Items []DefinitionNode
}

type DefinitionNode struct {
	Term []InlineNode
	Desc []InlineNode
}

type TableNode struct {
	Rows []TableRowNode
}

type TableCell struct {
	Inlines []InlineNode
}

type TableRowNode struct {
	Cells []TableCell
}

type PreformattedNode struct {
	Lines []string
}

type BlankLineNode struct{}
