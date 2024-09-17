package rapidreporter

import (
	"github.com/alecthomas/participle/v2"
)

type Report struct {
	Name         string       `parser:"'report' @String '{'"`
	Query        string       `parser:"'query' @String"`
	Params       []string     `parser:"'params' '[' @String* ']'"`
	RowProcesses []RowProcess `parser:"'rowprocesses' '[' @@* ']'"`
	Close        string       `parser:"'}'"`
}

type RowProcess struct {
	ColumnName string   `parser:"@String"`
	Columns    []string `parser:"@Ident*"`
	Operator   string   `parser:"@( '+' | '-' | '*' | '/' )"`
}

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
