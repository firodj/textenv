package db

import (
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	dbConn := New()
	fmt.Println(dbConn.URL)

	if !strings.HasPrefix(dbConn.URL, "sql://") {
		t.Error()
	}

	if !strings.HasSuffix(dbConn.URL, "/somedb") {
		t.Error()
	}
}
