package textenv

import "testing"

func Test_parseLine(t *testing.T) {
	for _, tt := range []struct {
		name   string
		l      string
		errmsg string
		setEnv func(string, string) error
		getEnv func(string) string
	}{
		{
			"when blank",
			"    ",
			"",
			nil,
			nil,
		},
		{
			"when comment",
			"   # this is a comments",
			"",
			nil,
			nil,
		},
		{
			"when no equal sign",
			"FORGET TO ADD EQ",
			"missing equal sign '=' on FORGET TO ADD EQ",
			nil,
			nil,
		},
		{
			"when key value",
			"WE = ARE=THE-BEST",
			"",
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
		{
			"when value has interpolate",
			"DB_USER=$USER",
			"",
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
	} {
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

			if tt.setEnv == nil {
				if countSetenv != 0 {
					t.Errorf("setenv should not be called")
				}
			} else {
				if countSetenv != 1 {
					t.Errorf("setenv should be called once")
				}
			}

			if tt.getEnv == nil {
				if countGetenv != 0 {
					t.Errorf("getenv should not be called")
				}
			} else {
				if countGetenv != 1 {
					t.Errorf("getenv should be called once")
				}
			}
		})
	}
}
