package textenv

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var GetEnv = os.Getenv
var SetEnv = os.Setenv

func GetProjectRootPath() (path string, err error) {
	_, filename, _, _ := runtime.Caller(0)
	path, err = filepath.Abs(filepath.Dir(filename) + "/..")
	return
}

func parseLine(l string) error {
	l = strings.TrimSpace(l)

	if len(l) == 0 {
		return nil
	}

	if strings.HasPrefix(l, "#") {
		return nil
	}

	s := strings.SplitN(l, "=", 2)
	if len(s) < 2 {
		return fmt.Errorf("missing equal sign '=' on %s", l)
	}

	k := strings.TrimSpace(s[0])
	v := strings.TrimSpace(s[1])

	if strings.HasPrefix(v, "$") {
		v = GetEnv(v[1:])
	}

	err := SetEnv(k, v)
	if err != nil {
		return fmt.Errorf("unable set env %s with value %s, error: %v", k, v, err)
	}
	return nil
}

func parseContents(b []byte) error {
	lines := strings.Split(string(b), "\n")

	for _, l := range lines {
		err := parseLine(l)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadTestEnv(textenv string) error {
	projPath, err := GetProjectRootPath()
	if err != nil {
		return fmt.Errorf("unable to get project root path: %v", err)
	}

	fullname := filepath.Join(projPath, textenv)
	b, err := os.ReadFile(fullname)
	if err != nil {
		return fmt.Errorf("unable to read %s, error: %v", fullname, err)
	}

	parseContents(b)
	return nil
}
