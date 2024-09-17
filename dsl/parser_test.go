package dsl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParser(t *testing.T) {
	parser, err := NewParser()
	if err != nil {
		t.Errorf("Error creating new parser: %v", err)
	}
	if parser == nil {
		t.Error("Parser is nil")
	}
}

func TestParse(t *testing.T) {
	parser, _ := NewParser()
	input := `report "customer report" {
query "select first_name, last_name from invoices where business_id = ?" 
params [
	"businessId"
	"userId"
	"startDate"
]
processes [
	"full_name" first_name last_name +
	"total" quantity price *
]
columns [
	label "First Name" key "first_name"
	label "Last Name" key "last_name"
]
}`
	report, err := parser.Parse(input)
	if err != nil {
		t.Errorf("Error parsing input: %v", err)
	}
	assert.Equal(t, `"customer report"`, report.Name)
	assert.Equal(t, `"select first_name, last_name from invoices where business_id = ?"`, report.Query)
	assert.Equal(t, `"businessId"`, report.Params[0])
	assert.Equal(t, `"userId"`, report.Params[1])
	assert.Equal(t, `"startDate"`, report.Params[2])
	assert.Equal(t, `"full_name"`, report.RowProcesses[0].ColumnName)
}
