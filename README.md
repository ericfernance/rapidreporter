This project is a basic report builder for Go.  It is inspired by KoolReport in PHP.

### Installation

To install the package, run:

```bash
go get github.com/ericfernance/rapidreporter
```


### Usage
Sample usage:

```go

dbCfg := mysql.Config{
    User:                 "someuser",
    Passwd:               "somesecret",
    Net:                  "tcp",
    Addr:                 "127.0.0.1:3306",
    DBName:               "adb",
}
db, err := sql.Open("mysql", dbCfg.FormatDSN())
if err != nil {
    panic(err)
}
rr, _ := rapidreporter.NewReporter(db)
rr.Query("SELECT id,first_name,last_name FROM clients where id = ?")
rr.Params([]interface{}{1})
rr.RowProcess(func(row map[string]interface{}) {
    row["full_name"] = fmt.Sprintf("%s %s", row["first_name"], row["last_name"])
})
rr.Columns([]rapidreporter.ReportColumn{
    {Label: "THE FIRST NAME", Key: "first_name"},
    {Label: "THE LAST NAME", Key: "last_name"},
    {Label: "THE FULL NAME", Key: "full_name"},
    {Label: "THE ID", Key: "id"},
})
rr.Run()
output := rr.Output("html")
fmt.Println(output)

```


Multiple calls to `RowProcess` can be changed with each call mutating the row data.

By default data is output as an html table.  Other output formats can be added by creating a function with the signature:

```go
func([]map[string]interface{}, []ReportColumn) string
```

the first parameter will be the data the second parameter will be all the column definitions.  The function should return a string representing the output.

The column definitions are expected as structs of type `ReportColumn` These are set using the `Columns` method.

```go
type ReportColumn struct {
    Label string
    Key   string
    ShowTotal bool
}
```

Where one or more columns have the `ShowTotal` flag set to true, the default visualiser output will include a total row at the end of the data.

It is the reponsibility of the visualiser to handle the `ShowTotal` flag.

A visualiser can be added to the reporter by calling `rr.Visualise(thefunction)`.  Only one visualiser can be added.



### Advanced Usage

To customize the visualization output, you can pass a custom visualizer function:

```go
reporter.Visualise(func(rows []map[string]interface{}, cols []ReportColumn) string {
    // Custom visualization logic here
    return "custom report"
})
```

This allows you to adapt the output format as per your requirements.

### Running Tests
To run the tests, use the following command:

```bash
go test ./...
```

Ensure you have sqlmock installed for proper test execution.