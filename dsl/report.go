package dsl

type Report struct {
	Name         string       `parser:"'report' @String '{'"`
	Query        string       `parser:"'query' @String"`
	Params       []string     `parser:"'params' '[' @String* ']'"`
	RowProcesses []RowProcess `parser:"'processes' '[' @@* ']'"`
	Columns      []Column     `parser:"'columns' '[' @@* ']'"`
	Close        string       `parser:"'}'"`
}

type RowProcess struct {
	ColumnName string   `parser:"@String"`
	Columns    []string `parser:"@Ident*"`
	Operator   string   `parser:"@( '+' | '-' | '*' | '/' )"`
}

type Column struct {
	Label string `parser:"'label' @String"`
	Key   string `parser:"'key' @String"`
}
