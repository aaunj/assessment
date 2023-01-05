package expense

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUnitUpdateById(t *testing.T) {
	t.Run("Update-Expense-by-ID", func(t *testing.T) {
		e := echo.New()
		reqBody := `{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`

		req, err := http.NewRequest(http.MethodPut, "/expenses/1", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when making API request", err)
		}
		//req.Header.Set(echo.HeaderAuthorization, "November 10, 2009")
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		c.SetPath("/expenses/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		conn, mock, err := sqlmock.New()
		db = conn
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(1, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"}))

		mock.ExpectQuery("UPDATE expenses SET (.+) WHERE (.+)").WithArgs("apple smoothie", float64(89), "no discount", pq.Array([]string{"beverage"}), 1).WillReturnRows(rows)

		if err = UpdateByIdHandler(c); err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		var expen Expense
		json.NewDecoder(resp.Body).Decode(&expen)
		assert.Equal(t, 1, expen.ID)
		assert.Equal(t, "apple smoothie", expen.Title)
		assert.Equal(t, float64(89), expen.Amount)
		assert.Equal(t, "no discount", expen.Note)
		assert.Equal(t, []string{"beverage"}, expen.Tags)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})
}
