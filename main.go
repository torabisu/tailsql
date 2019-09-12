package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
)

// clearScreeen simply clears the terminal and returns the cursor to the top
func clearScreen() {
	log.Println("Clearing the screen...")
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// getColumns counts the number of columns returned in the query, and returns the column headers
// and an area in which you can store the values
func getColumns(rows *sql.Rows) ([]string, []interface{}) {
	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	vals := make([]interface{}, len(cols))
	for i, _ := range cols {
		vals[i] = new(sql.RawBytes)
	}
	return cols, vals
}

// processQuery runs the query, builds up the returned dataset and displays the result
func processQuery(db *sql.DB, query *string) {
	rows, err := db.Query(*query)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	cols, vals := getColumns(rows)
	dataset := [][]string{}

	for rows.Next() {
		row := []string{}
		rows.Scan(vals...)
		for _, col := range vals {
			row = append(row, string(*col.(*sql.RawBytes)))
		}
		dataset = append(dataset, row)
	}
	renderTable(cols, dataset)
}

// renderTable given column headers and a dataset, renders an ASCII table
func renderTable(cols []string, dataset [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(cols)
	table.AppendBulk(dataset)
	table.Render()
}

func main() {
	var username = flag.String("username", "", "Database username")
	var password = flag.String("password", "", "Database password")
	var database = flag.String("database", "", "Database name")
	var query = flag.String("query", "", "The SQL query to run against the database")
	var clear = flag.Bool("clear", false, "Set to true if you want to clear the terminal between runs")
	var sleep = flag.Int("sleep", 3, "Number of seconds to sleep for before rerunning the query")

	flag.Parse()
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", *username, *password, *database))

	if err != nil {
		log.Fatal(err)
	}

	for {
		if *clear {
			clearScreen()
		}
		processQuery(db, query)
		fmt.Printf("Querying again in %v seconds...\n", *sleep)
		time.Sleep(time.Second * time.Duration(*sleep))
	}
}
