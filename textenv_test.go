package textenv

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetScriptPath(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	expected, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		t.Error(err)
	}

	actual, err := GetScriptPath(0)
	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Errorf("unexpected %s with %s", actual, expected)
	}
}

type TestCase interface {
	Run(t *testing.T)
}

type testCase struct {
	name   string
	l      string
	setEnv func(string, string) error
	getEnv func(string) string
}

func (tt *testCase) runTest(t *testing.T, tc TestCase) {
	t.Run(tt.name, func(t *testing.T) {
		oldSetenv := SetEnv
		oldGetenv := GetEnv
		defer func() {
			SetEnv = oldSetenv
			GetEnv = oldGetenv
		}()
		countSetenv := 0
		countGetenv := 0

		SetEnv = func(k string, v string) error {
			countSetenv += 1
			if tt.setEnv != nil {
				return tt.setEnv(k, v)
			}
			return nil
		}
		GetEnv = func(k string) (v string) {
			countGetenv += 1
			if tt.getEnv != nil {
				return tt.getEnv(k)
			}
			return
		}

		tc.Run(t)

		if tt.setEnv == nil {
			if countSetenv > 0 {
				t.Errorf("setenv should not be called")
			}
		} else {
			if countSetenv == 0 {
				t.Errorf("setenv should be called")
			}
		}

		if tt.getEnv == nil {
			if countGetenv > 0 {
				t.Errorf("getenv should not be called")
			}
		} else {
			if countGetenv == 0 {
				t.Errorf("getenv should be called")
			}
		}
	})
}

type parseLineTestCase struct {
	testCase
	errmsg string
}

func (tt parseLineTestCase) Run(t *testing.T) {
	err := parseLine(tt.l)
	if err != nil {
		if tt.errmsg != err.Error() {
			t.Errorf("unexpected error: %s, with: %s", err.Error(), tt.errmsg)
		}
	} else {
		if tt.errmsg != "" {
			t.Errorf("unmeet expected error: %s", tt.errmsg)
		}
	}
}

func Test_parseLine(t *testing.T) {
	for _, tt := range []parseLineTestCase{
		{
			testCase{
				"when blank",
				"    ",
				nil,
				nil,
			},
			"",
		},
		{
			testCase{
				"when comment",
				"   # this is a comments",
				nil,
				nil,
			},
			"",
		},
		{
			testCase{
				"when no equal sign",
				"FORGET TO ADD EQ",
				nil,
				nil,
			},
			"missing equal sign '=' on FORGET TO ADD EQ",
		},
		{
			testCase{
				"when key value",
				"WE = ARE=THE-BEST",
				func(k string, v string) error {
					if k != "WE" {
						t.Errorf("unexpected k: %s", k)
					}
					if v != "ARE=THE-BEST" {
						t.Errorf("unexpected v: %s", v)
					}
					return nil
				},
				nil,
			},
			"",
		},
		{
			testCase{
				"when value has interpolate",
				"DB_USER=$(USER)",
				func(k string, v string) error {
					if k != "DB_USER" {
						t.Errorf("unexpected k: %s", k)
					}
					if v != "cloudsql" {
						t.Errorf("unexpected v: %s", v)
					}
					return nil
				},
				func(k string) (v string) {
					if k != "USER" {
						t.Errorf("unexpected k: %s", k)
					}
					v = "cloudsql"
					return
				},
			},
			"",
		},
	} {
		tt.runTest(t, tt)
	}
}

type subReplaceTestCase struct {
	testCase
	result string
}

func (tt subReplaceTestCase) Run(t *testing.T) {
	result := subReplace(tt.l)
	if result != tt.result {
		t.Errorf("unexpected result: %s, with: %s", result, tt.result)
	}
}

func Test_subReplace(t *testing.T) {
	for _, tt := range []subReplaceTestCase{
		{
			testCase{
				"when blank",
				"",
				nil,
				nil,
			},
			"",
		},
		{
			testCase{
				"when there is",
				"AKU$(ADALAH)DIA$$(TAPIBUKAN)DAN$(LAGI)",
				nil,
				func(k string) string {
					switch k {
					case "ADALAH":
						return "SESUATU"
					case "LAGI":
						return "SUDAH"
					default:
						t.Errorf("unexpected getenv(%s)", k)
					}
					return ""
				},
			},
			"AKUSESUATUDIA$(TAPIBUKAN)DANSUDAH",
		},
	} {
		tt.testCase.runTest(t, tt)
	}
}
