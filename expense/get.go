package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetByIdHandler(c echo.Context) error {
	var expen Expense
	id, _ := strconv.Atoi(c.Param("id"))

	row := db.QueryRow("SELECT * FROM expenses WHERE id=$1", id)
	err := row.Scan(&expen.ID, &expen.Title, &expen.Amount, &expen.Note, pq.Array(&expen.Tags))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	//fmt.Printf("data : % #v\n", expen)

	return c.JSON(http.StatusOK, expen)

}

func GetAllHandler(c echo.Context) error {
	var expenses []Expense

	rows, err := db.Query("SELECT * FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	defer rows.Close()
	for rows.Next() {
		var expen Expense
		err := rows.Scan(&expen.ID, &expen.Title, &expen.Amount, &expen.Note, pq.Array(&expen.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		}
		expenses = append(expenses, expen)
	}

	//fmt.Printf("all : % #v\n", expenses)

	return c.JSON(http.StatusOK, expenses)

}
