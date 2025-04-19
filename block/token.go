package block

type TokenType int

const (
	TokenBlank TokenType = iota
	TokenOrderedList
	TokenUnOrderedList
	TokenDefinitionList
	TokenHeading
	TokenTablePipeRow
	TokenTableCommaRow
	TokenPreformatted
	TokenParagraph
	TokenMaybeDirective
	TokenEOF
)

type Token struct {
	Type    TokenType
	Content string
	Level   int
}

func newToken(typ TokenType) Token {
	return Token{
		Type:    typ,
		Content: "",
		Level:   0,
	}
}

func newTokenWithContent(typ TokenType, content string) Token {
	return Token{
		Type:    typ,
		Content: content,
		Level:   0,
	}
}

func newTokenWithLevel(typ TokenType, content string, level int) Token {
	return Token{
		Type:    typ,
		Content: content,
		Level:   level,
	}
}

func (t TokenType) String() string {
	switch t {
	case TokenBlank:
		return "Blank"
	case TokenOrderedList:
		return "OrderedList"
	case TokenUnOrderedList:
		return "UnOrderedList"
	case TokenDefinitionList:
		return "DefinitionList"
	case TokenTablePipeRow:
		return "TablePipeRow"
	case TokenTableCommaRow:
		return "TableCommaRow"
	case TokenHeading:
		return "Heading"
	case TokenParagraph:
		return "Paragraph"
	case TokenPreformatted:
		return "Preformatted"
	case TokenMaybeDirective:
		return "MaybeDirective"
	case TokenEOF:
		return "EOF"
	default:
		return "Unknown"
	}
}
