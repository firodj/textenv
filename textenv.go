package textenv

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var GetEnv = os.Getenv
var SetEnv = os.Setenv

func GetScriptPath(skip int) (path string, err error) {
	_, filename, _, _ := runtime.Caller(skip+1)
	path, err = filepath.Abs(filepath.Dir(filename))
	return
}

func subReplace(l string) string {
	re := regexp.MustCompile(`(^|.)\$\(([A-Za-z0-9_]+)\)`)
	repl := func(groups []string) string {
		if groups[1] == "$" {
			return "$(" +  groups[2] + ")";
		}
		return groups[1] + GetEnv(groups[2])
	}
	result := ""

	g0 := 0
	for _, g := range re.FindAllStringSubmatchIndex(l, -1) {
		groups := []string{}
		for i := 0; i < len(g); i += 2 {
			groups = append(groups, l[g[i]:g[i+1]])
		}

		result += l[g0:g[0]] + repl(groups)
		g0 = g[1]
	}
	result += l[g0:]

	return result
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
	v := subReplace(strings.TrimSpace(s[1]))

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
	projPath, err := GetScriptPath(1)
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
