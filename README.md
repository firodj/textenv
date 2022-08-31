# textenv
golang load text env to env variable, allow interpolate, specifically for run testing thru IDE

# test

```
$ go test -v -cover
```

# development
```
$ go fmt
```

# example

File `PROJECT_DIR/test.env`, contents:

```env
# testing environment
DBUSER=$(USER)
DBHOST=localhost
DBURL=sql://$(DBUSER)@$(DBHOST)/somedb
```

File `PROJECT_DIR/db/init_test.go`, contents:

```go
func init() {
  oriSetenv := textenv.SetEnv
  defer func() {
    textenv.SetEnv = oriSetenv
  }()

  textenv.SetEnv = func(k string, v string) error {
    fmt.Printf("env %s=%s\n", k, v)
    return oriSetenv(k, v)
  }

  err := textenv.LoadTestEnv("../test.env")
  if err != nil {
    if !strings.Contains(err.Error(), "no such file") {
      panic(err)
    } else {
      fmt.Println("skip load test.env")
    }
  }
}
```
