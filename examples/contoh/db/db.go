package db

import "os"

type DBConn struct {
	URL string
}

func New() *DBConn {
	return &DBConn{
		URL: os.Getenv("DBURL"),
	}
}
