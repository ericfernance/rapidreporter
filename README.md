This project is a basic report builder for Go.  It is inspired by KoolReport in the PHP world.

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
rr.Process(func(row map[string]interface{}) {
    row["full_name"] = fmt.Sprintf("%s %s", row["first_name"], row["last_name"])
})
rr.Run()
output := rr.Output("html")
fmt.Println(output)

```


Multiple calls to `Process` can be changed with each call mutating the row data.

By default data is output as an html table.  Other output formats can be added by creating a function with the signature:

```go
func([]map[string]interface{}, map[string]string) string
```

the first parameter will be the data the second paramater will be all the column definitions.  The function should return a string representing the oputput.

These can be added to the reporter by calling `rr.Visualise(thefunction)`



