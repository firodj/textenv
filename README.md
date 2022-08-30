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

```go
func init() {

  oldSetenv := textenv.SetEnv

  defer func() {
		textenv.SetEnv = oldSetenv
	}()

  textenv.SetEnv = func(k string, v string) error {
		fmt.Printf("env %s=%s\n", k, v)
		return oldSetenv(k, v)
	}

	err := textenv.LoadTestEnv("../test.env")
	if err != nil {
		panic(err)
	}
}
```
