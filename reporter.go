package rapidreporter

import (
	"database/sql"
	"fmt"
)

type Reporter struct {
	connection   *sql.DB
	query        string
	params       []interface{}
	rowProcesses []func(map[string]interface{})
	rows         []map[string]interface{}
	columnDefs   []ReportColumn
	visualise    func([]map[string]interface{}, []ReportColumn) string
}
type ReportColumn struct {
	Label     string
	Key       string
	ShowTotal bool
}

func tableVisualiser(rows []map[string]interface{}, columnDefs []ReportColumn) string {
	html := "<table><thead><tr>"
	showFooter := false
	totals := make(map[string]float64)
	for col := range columnDefs {
		label := columnDefs[col].Label
		html += "<th>" + label + "</th>"
		if columnDefs[col].ShowTotal {
			showFooter = true
			totals[label] = 0
		}
	}
	html += "</tr></thead><tbody>"
	for _, row := range rows {
		html += "<tr>"
		for col := range columnDefs {
			key := columnDefs[col].Key
			value := row[key]
			html += fmt.Sprintf("<td>%v</td>", value)
			if columnDefs[col].ShowTotal {
				switch v := value.(type) {
				case int64:
					totals[columnDefs[col].Label] += float64(v)
				case float64:
					totals[columnDefs[col].Label] += v
				default:
					totals[columnDefs[col].Label] += 0
				}
			}
		}
		html += "</tr>"
	}
	html += "</tbody>"

	// Add the footer row
	if showFooter {
		html += "<tfoot><tr>"
		for col := range columnDefs {
			label := columnDefs[col].Label
			if columnDefs[col].ShowTotal {
				html += fmt.Sprintf("<td>%v</td>", totals[label])
			} else {
				html += "<td></td>"
			}
		}
		html += "</tr></tfoot>"
	}

	html += "</table>"
	return html
}

func NewReporter(db *sql.DB) (*Reporter, error) {
	return &Reporter{
		connection: db,
		visualise:  tableVisualiser,
	}, nil
}

func (r *Reporter) Visualise(visualise func([]map[string]interface{}, []ReportColumn) string) *Reporter {
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
		for _, process := range r.rowProcesses {
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

func (r *Reporter) RowProcess(process func(map[string]interface{})) *Reporter {
	r.rowProcesses = append(r.rowProcesses, process)
	return r
}

func (r *Reporter) Output(format string) string {
	return r.visualise(r.rows, r.columnDefs)
}

func (r *Reporter) Columns(columnDefs []ReportColumn) *Reporter {
	r.columnDefs = columnDefs
	return r
}
