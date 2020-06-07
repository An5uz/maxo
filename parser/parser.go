package parser

import (
	"fmt"
	"github.com/an5uz/maxo/lexer"
	"strings"

)

func Parse(text string) (string, error) {
	p := parser{
		lex: lexer.Lex(text),
	}

	p.parse()

	if p.errItem != nil {
		return "", fmt.Errorf("error processing the following %q", p.errItem.Value)
	}

	return p.result, nil
}

type parser struct {
	result  string
	lex     *lexer.Lexer
	errItem *lexer.Item
}

func (p *parser) parse() {
	sb := strings.Builder{}

	for item := range p.lex.Items {
		switch item.Kind {
		case lexer.ItemEOF:
			p.result = sb.String()
			return
		case lexer.ItemError:
			p.errItem = &item
			return
		case lexer.ItemText:
			sb.WriteString(lexer.Reverse(item.Value))

		case lexer.ItemWhiteSpace:
			sb.WriteString(item.Value)
		}
	}
}
