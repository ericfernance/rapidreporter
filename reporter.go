package rapidreporter

import (
	"database/sql"
	"fmt"
)

type Reporter struct {
	connection *sql.DB
	query      string
	params     []interface{}
	processes  []func(map[string]interface{})
	rows       []map[string]interface{}
	columnDefs []map[string]interface{}
	visualise  func([]map[string]interface{}, []map[string]interface{}) string
}

func tableVisualiser(rows []map[string]interface{}, columnDefs []map[string]interface{}) string {
	html := "<table><thead><tr>"
	for col := range columnDefs {
		label := columnDefs[col]["label"].(string)
		html += "<th>" + label + "</th>"
	}
	html += "</tr></thead><tbody>"
	for _, row := range rows {
		html += "<tr>"
		for col := range columnDefs {
			key := columnDefs[col]["key"].(string)
			value := row[key]
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

func (r *Reporter) Visualise(visualise func([]map[string]interface{}, []map[string]interface{}) string) *Reporter {
	r.visualise = visualise
	return r
}

func (r *Reporter) Params(params []interface{}) *Reporter {
	r.params = params
	return r
}

func (r *Reporter) Run() (*Reporter, error) {
	rows, err := r.connection.Query(r.query, r.params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	valuePtrs := make([]interface{}, len(columns))

	// Process each row
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range columns {
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
	return r.visualise(r.rows, r.columnDefs)
}

func (r *Reporter) Columns(columnDefs []map[string]interface{}) *Reporter {
	r.columnDefs = columnDefs
	return r
}
