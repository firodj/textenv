package main

import (
	"contoh/db"
	"fmt"
)

func main() {
	contoh()
}

func contoh() {
	dbConn := db.New()
	fmt.Printf("Contoh app, DB URL = %v\n", dbConn.URL)
}
