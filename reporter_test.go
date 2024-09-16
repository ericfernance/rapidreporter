package rapidreporter

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewReporter(t *testing.T) {
	// Create a mock database connection
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a new reporter
	reporter, err := NewReporter(db)
	assert.NoError(t, err)
	assert.NotNil(t, reporter)
	assert.Equal(t, db, reporter.connection)
	assert.NotNil(t, reporter.visualise)
}

func TestReporter_Query(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	reporter, err := NewReporter(db)
	assert.NoError(t, err)

	// Set the query
	query := "SELECT * FROM test_table"
	reporter.Query(query)

	assert.Equal(t, query, reporter.query)
}

func TestReporter_Params(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	reporter, err := NewReporter(db)
	assert.NoError(t, err)

	params := []interface{}{"param1", "param2"}
	reporter.Params(params)

	assert.Equal(t, params, reporter.params)
}

func TestReporter_Run(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Mock the database response
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "John Doe").
		AddRow(2, "Jane Doe")
	mock.ExpectQuery("SELECT \\* FROM test_table").WillReturnRows(rows)

	// Create a new reporter and run the query
	reporter, err := NewReporter(db)
	assert.NoError(t, err)

	reporter.Query("SELECT * FROM test_table")
	reporter, err = reporter.Run()
	assert.NoError(t, err)

	// Check if rows are processed correctly
	expectedRows := []map[string]interface{}{
		{"id": int64(1), "name": "John Doe"},
		{"id": int64(2), "name": "Jane Doe"},
	}
	assert.Equal(t, expectedRows, reporter.rows)

	// Ensure that all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReporter_Visualise(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	reporter, err := NewReporter(db)
	assert.NoError(t, err)

	// Custom visualisation function for testing
	customVisualise := func(rows []map[string]interface{}, columnDefs []ReportColumn) string {
		return "custom visualisation"
	}

	reporter.Visualise(customVisualise)

	assert.Equal(t, "custom visualisation", reporter.Output("custom"))
}

func TestReporter_Columns(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	reporter, err := NewReporter(db)
	assert.NoError(t, err)

	columnDefs := []ReportColumn{
		{Label: "ID", Key: "id", ShowTotal: false},
		{Label: "Name", Key: "name", ShowTotal: false},
	}

	reporter.Columns(columnDefs)

	assert.Equal(t, columnDefs, reporter.columnDefs)
}

func TestTableVisualiser(t *testing.T) {
	columnDefs := []ReportColumn{
		{Label: "ID", Key: "id", ShowTotal: false},
		{Label: "Amount", Key: "amount", ShowTotal: true},
	}

	rows := []map[string]interface{}{
		{"id": int64(1), "amount": float64(100)},
		{"id": int64(2), "amount": float64(150)},
	}

	expectedOutput := `<table><thead><tr><th>ID</th><th>Amount</th></tr></thead><tbody><tr><td>1</td><td>100</td></tr><tr><td>2</td><td>150</td></tr></tbody><tfoot><tr><td></td><td>250</td></tr></tfoot></table>`

	result := tableVisualiser(rows, columnDefs)
	assert.Equal(t, expectedOutput, result)
}

func TestTableVisualiser_MixedTypes(t *testing.T) {
	columnDefs := []ReportColumn{
		{Label: "ID", Key: "id", ShowTotal: false},
		{Label: "Amount", Key: "amount", ShowTotal: true},
		{Label: "Name", Key: "name", ShowTotal: false},
	}

	rows := []map[string]interface{}{
		{"id": int64(1), "amount": float64(100.5), "name": "John"},
		{"id": int64(2), "amount": int64(50), "name": "Jane"},
		{"id": int64(3), "amount": float64(200.25), "name": "Doe"},
	}

	expectedOutput := `<table><thead><tr><th>ID</th><th>Amount</th><th>Name</th></tr></thead><tbody><tr><td>1</td><td>100.5</td><td>John</td></tr><tr><td>2</td><td>50</td><td>Jane</td></tr><tr><td>3</td><td>200.25</td><td>Doe</td></tr></tbody><tfoot><tr><td></td><td>350.75</td><td></td></tr></tfoot></table>`

	result := tableVisualiser(rows, columnDefs)
	assert.Equal(t, expectedOutput, result)
}

func TestTableVisualiser_IncorrectAmountString(t *testing.T) {
	columnDefs := []ReportColumn{
		{Label: "ID", Key: "id", ShowTotal: false},
		{Label: "Amount", Key: "amount", ShowTotal: true},
		{Label: "Name", Key: "name", ShowTotal: false},
	}

	rows := []map[string]interface{}{
		{"id": int64(1), "amount": float64(100.5), "name": "John"},
		{"id": int64(2), "amount": "invalid", "name": "Jane"}, // Incorrect amount as string
		{"id": int64(3), "amount": float64(200.25), "name": "Doe"},
	}

	expectedOutput := `<table><thead><tr><th>ID</th><th>Amount</th><th>Name</th></tr></thead><tbody><tr><td>1</td><td>100.5</td><td>John</td></tr><tr><td>2</td><td>invalid</td><td>Jane</td></tr><tr><td>3</td><td>200.25</td><td>Doe</td></tr></tbody><tfoot><tr><td></td><td>300.75</td><td></td></tr></tfoot></table>`

	result := tableVisualiser(rows, columnDefs)
	assert.Equal(t, expectedOutput, result)
}
