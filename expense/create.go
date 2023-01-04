package expense

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func CreateHandler(c echo.Context) error {
	var expen Expense
	err := c.Bind(&expen)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	err = queryRow(db, &expen)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	fmt.Printf("id : % #v\n", expen)

	return c.JSON(http.StatusCreated, expen)
}

func queryRow(db *sql.DB, expen *Expense) error {
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id;", expen.Title, expen.Amount, expen.Note, pq.Array(expen.Tags))
	return row.Scan(&expen.ID)
}
