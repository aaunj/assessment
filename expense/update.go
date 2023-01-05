package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func UpdateByIdHandler(c echo.Context) error {
	var expen Expense

	err := c.Bind(&expen)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	expen.ID, _ = strconv.Atoi(c.Param("id"))

	row := db.QueryRow("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id,title,amount,note,tags;",
		expen.Title, expen.Amount, expen.Note, pq.Array(expen.Tags), expen.ID)
	err = row.Scan(&expen.ID, &expen.Title, &expen.Amount, &expen.Note, pq.Array(&expen.Tags))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	//fmt.Printf("updateById : % #v\n", expen)

	return c.JSON(http.StatusOK, expen)

}
