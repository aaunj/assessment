package expense

import (
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

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES($1, $2, $3, $4) RETURNING id;",
		expen.Title, expen.Amount, expen.Note, pq.Array(expen.Tags))
	err = row.Scan(&expen.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	//fmt.Printf("Create : % #v\n", expen)

	return c.JSON(http.StatusCreated, expen)
}
