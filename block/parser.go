package block

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/miso24/puki2md/inline"
	"github.com/miso24/puki2md/model"
)

type Parser struct {
	t   *Tokenizer
	buf *Token
}

func (p *Parser) peekToken() Token {
	if p.buf == nil {
		tok := p.t.TokenizeLine()
		p.buf = &tok
	}
	return *p.buf
}

func (p *Parser) readToken() Token {
	if p.buf != nil {
		tok := *p.buf
		p.buf = nil
		return tok
	}
	return p.t.TokenizeLine()
}

func (p *Parser) parseParagraph() model.BlockNode {
	tok := p.readToken()
	inlines := inline.Parse(tok.Content)
	return model.ParagraphNode{Inlines: inlines}
}

func (p *Parser) parseHeading() model.BlockNode {
	tok := p.readToken()
	inlines := inline.Parse(tok.Content)
	return model.HeadingNode{Level: tok.Level, Inlines: inlines}
}

func (p *Parser) parseList(ordered bool) model.BlockNode {
	tok := p.readToken()
	inlines := inline.Parse(tok.Content)
	return model.ListNode{Level: tok.Level, Inlines: inlines, Ordered: ordered}
}

func (p *Parser) parseDefinition() model.BlockNode {
	items := []model.DefinitionNode{}
	for {
		tok := p.peekToken()
		if tok.Type != TokenDefinitionList {
			break
		}

		elems := strings.SplitN(p.readToken().Content[1:], "|", 2)
		term := inline.Parse(elems[0])
		desc := inline.Parse(elems[1])
		items = append(items, model.DefinitionNode{Term: term, Desc: desc})
	}
	return model.DefinitionListNode{Items: items}
}

func (p *Parser) parsePreformatted() model.BlockNode {
	lines := []string{}
	for {
		tok := p.peekToken()
		if tok.Type != TokenPreformatted {
			break
		}
		lines = append(lines, p.readToken().Content)
	}
	return model.PreformattedNode{Lines: lines}
}

func (p *Parser) parsePipeTable() model.BlockNode {
	rows := []model.TableRowNode{}
	for {
		tok := p.peekToken()
		if tok.Type != TokenTablePipeRow {
			break
		}

		cells := []model.TableCell{}
		row := strings.Trim(tok.Content, "|")
		for _, cell := range strings.Split(row, "|") {
			inlines := inline.Parse(cell)
			cells = append(cells, model.TableCell{Inlines: inlines})
		}
		rows = append(rows, model.TableRowNode{Cells: cells})
		p.readToken()
	}
	return model.TableNode{Rows: rows}
}

func (p *Parser) parseCommaTable() model.BlockNode {
	var s strings.Builder
	for {
		tok := p.peekToken()
		if tok.Type != TokenTableCommaRow {
			break
		}
		s.WriteString(p.readToken().Content[1:] + "\n")
	}

	r := csv.NewReader(strings.NewReader(s.String()))
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	rows := []model.TableRowNode{}
	for _, record := range records {
		cells := []model.TableCell{}
		for _, cell := range record {
			inlines := inline.Parse(cell)
			cells = append(cells, model.TableCell{Inlines: inlines})
		}
		rows = append(rows, model.TableRowNode{Cells: cells})
	}
	return model.TableNode{Rows: rows}
}

func (p *Parser) Parse() []model.BlockNode {
	nodes := []model.BlockNode{}
	for {
		tok := p.peekToken()
		if tok.Type == TokenEOF {
			break
		}
		fmt.Printf("[%s] %s\n", tok.Type, tok.Content)

		switch tok.Type {
		case TokenParagraph:
			nodes = append(nodes, p.parseParagraph())
		case TokenHeading:
			nodes = append(nodes, p.parseHeading())
		case TokenUnOrderedList:
			nodes = append(nodes, p.parseList(false))
		case TokenOrderedList:
			nodes = append(nodes, p.parseList(true))
		case TokenDefinitionList:
			nodes = append(nodes, p.parseDefinition())
		case TokenPreformatted:
			nodes = append(nodes, p.parsePreformatted())
		case TokenTablePipeRow:
			nodes = append(nodes, p.parsePipeTable())
		case TokenTableCommaRow:
			nodes = append(nodes, p.parseCommaTable())
        case TokenMaybeDirective:
            // TODO: implement
		case TokenBlank:
			p.readToken()
			nodes = append(nodes, model.BlankLineNode{})
		default:
			p.readToken()
		}
	}
	return nodes
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		t: NewTokenizer(r),
	}
}
