package expense

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetByIdHandler(c echo.Context) error {
	var expen Expense
	expen.ID, _ = strconv.Atoi(c.Param("id"))

	err := getByID(db, &expen)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	fmt.Printf("data : % #v\n", expen)

	return c.JSON(http.StatusOK, expen)

}

func getByID(db *sql.DB, expen *Expense) error {
	row := db.QueryRow("SELECT * FROM expenses WHERE id=$1", expen.ID)
	return row.Scan(&expen.ID, &expen.Title, &expen.Amount, &expen.Note, pq.Array(&expen.Tags))
}
