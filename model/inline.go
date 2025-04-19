package model

type InlineNode interface{}

type Strong struct {
	Inlines []InlineNode
}

type Italic struct {
	Inlines []InlineNode
}

type PageLink struct {
	Name   string
	Target string
}

type FootNote struct {
	Inlines []InlineNode
}

type PlainText struct {
	Content string
}

type Attach struct {
	FileName string
}

type RawCode struct {
	Code string
}
