package textenv_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/firodj/textenv"
)

func Example() {
  os.Setenv("CUSTOM_USER", "anonymous")

	oriSetenv := textenv.SetEnv
  defer func() {
    textenv.SetEnv = oriSetenv
  }()

  textenv.SetEnv = func(k string, v string) error {
    fmt.Printf("env %s=%s\n", k, v)
    return oriSetenv(k, v)
  }

  err := textenv.LoadTestEnv("./fixture/test.env")
  if err != nil {
    if !strings.Contains(err.Error(), "no such file") {
      panic(err)
    } else {
      fmt.Println("skip load test.env")
    }
  }

  err = textenv.LoadTestEnv("./test_notexist.env")
  if err != nil {
    if !strings.Contains(err.Error(), "no such file") {
      panic(err)
    } else {
      fmt.Println("skip load test_notexist.env")
    }
  }

  /* Output:
env DBUSER=anonymous
env DBHOST=localhost
env DBURL=sql://anonymous@localhost/somedb
skip load test_notexist.env
*/
}

