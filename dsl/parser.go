package dsl

import "github.com/alecthomas/participle/v2"

type Parser struct {
	parser *participle.Parser[Report]
}

func NewParser() (*Parser, error) {
	parser, err := participle.Build[Report]()
	if err != nil {
		return nil, err
	}
	return &Parser{
		parser: parser,
	}, nil
}

func (p *Parser) Parse(input string) (*Report, error) {
	report, err := p.parser.ParseString("", input)
	if err != nil {
		return nil, err
	}
	return report, nil
}
