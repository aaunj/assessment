package expense

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func InitDB() {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	db = conn
	defer db.Close()
	createExpensesTable()
}

func createExpensesTable() {
	createExpensesTable := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`

	if _, err := db.Exec(createExpensesTable); err != nil {
		log.Fatal("can't create table ", err)
	}
}
