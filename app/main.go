package main

import (
	"fmt"
	"log"
	"os"
	// Available if you need it!
	// "github.com/xwb1989/sqlparser"
)

// Usage: your_program.sh sample.db .dbinfo
func main() {
	databaseFilePath := os.Args[1]
	command := os.Args[2]

	switch command {
	case ".dbinfo":
		databaseFile, err := os.Open(databaseFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// You can use print statements as follows for debugging, they'll be visible when running tests.
		// fmt.Println("Logs from your program will appear here!")

		dbHeader, err := ParseDatabaseHeader(databaseFile)
		if err != nil {
			log.Fatalf("Failed to read the database header %s", err)
		}
		page := ParsePageHeader(databaseFile)

		// Uncomment this to pass the first stage
		fmt.Printf("database page size: %v\n", dbHeader.PageSize)
		fmt.Printf("Page type is: %s\n", page.Type.ToString())
		fmt.Printf("Cell nums is: %d\n", page.CellNums)
	default:
		fmt.Println("Unknown command", command)
		os.Exit(1)
	}
}
