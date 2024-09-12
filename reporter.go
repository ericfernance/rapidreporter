package rapidreporter

import (
	"database/sql"
	"fmt"
)

type Reporter struct {
	connection  *sql.DB
	query       string
	params      []interface{}
	processes   []func(map[string]interface{})
	rows        []map[string]interface{}
	columnTypes map[string]string
	visualise   func([]map[string]interface{}, map[string]string) string
}

func tableVisualiser(rows []map[string]interface{}, columnTypes map[string]string) string {
	html := "<table><thead><tr>"
	for key := range rows[0] {
		html += "<th>" + key + "</th>"
	}
	html += "</tr></thead><tbody>"
	for _, row := range rows {
		html += "<tr>"
		for _, value := range row {
			html += fmt.Sprintf("<td>%v</td>", value)
		}
		html += "</tr>"
	}
	html += "</tbody></table>"
	return html
}

func NewReporter(db *sql.DB) (*Reporter, error) {
	return &Reporter{
		connection: db,
		visualise:  tableVisualiser,
	}, nil
}

func (r *Reporter) Visualise(visualise func([]map[string]interface{}, map[string]string) string) *Reporter {
	r.visualise = visualise
	return r
}

func (r *Reporter) Params(params []interface{}) *Reporter {
	r.params = params
	return r
}

func (r *Reporter) Run() (*Reporter, error) {
	r.columnTypes = make(map[string]string)
	rows, err := r.connection.Query(r.query, r.params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	valuePtrs := make([]interface{}, len(columns))

	// Process each row
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range columns {
			r.columnTypes[columns[i]] = columnTypes[i].DatabaseTypeName()
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		dataRow := make(map[string]interface{})

		for i, col := range columns {
			val := values[i]

			if b, ok := val.([]byte); ok {
				dataRow[col] = string(b)
			} else {
				dataRow[col] = val
			}
		}

		// Process the rows through any user defined processes
		for _, process := range r.processes {
			process(dataRow)
		}

		r.rows = append(r.rows, dataRow)
	}
	return r, nil
}

func (r *Reporter) Query(query string) *Reporter {
	r.query = query
	return r
}

func (r *Reporter) Process(process func(map[string]interface{})) *Reporter {
	r.processes = append(r.processes, process)
	return r
}

func (r *Reporter) Output(format string) string {
	return r.visualise(r.rows, r.columnTypes)
}
