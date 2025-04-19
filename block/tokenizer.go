package block

import (
	"bufio"
	"io"
	"strings"
)

type Tokenizer struct {
	s *bufio.Scanner
}

func tokenizeListOrHeading(line string, prefix byte) (int, string) {
	level := 0
	for i := 0; i < len(line); i++ {
		if line[i] != prefix {
			break
		}
		level++
	}

	return level, strings.TrimSpace(line[level:])
}

func (t *Tokenizer) TokenizeLine() Token {
	if !t.s.Scan() {
		return newToken(TokenEOF)
	}

	line := t.s.Text()
	line = strings.TrimRight(line, " ")
	if len(line) == 0 {
		return newToken(TokenBlank)
	}

	switch line[0] {
	case '*':
		level, content := tokenizeListOrHeading(line, '*')
		return newTokenWithLevel(TokenHeading, content, level)
	case '-':
		level, content := tokenizeListOrHeading(line, '-')
		return newTokenWithLevel(TokenUnOrderedList, content, level)
	case '+':
		level, content := tokenizeListOrHeading(line, '+')
		return newTokenWithLevel(TokenOrderedList, content, level)
	case ':':
		return newTokenWithContent(TokenDefinitionList, line)
	case ' ':
		return newTokenWithContent(TokenPreformatted, line[1:])
	case '|':
		return newTokenWithContent(TokenTablePipeRow, line)
	case ',':
		return newTokenWithContent(TokenTableCommaRow, line)
	case '#':
		return newTokenWithContent(TokenMaybeDirective, line)
	}
	return newTokenWithContent(TokenParagraph, line)
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		s: bufio.NewScanner(r),
	}
}
