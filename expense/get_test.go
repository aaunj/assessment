package expense

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetById(t *testing.T) {
	t.Run("Get-Expense-by-ID", func(t *testing.T) {
		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/expenses/1", nil)
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
			AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"}))

		mock.ExpectQuery("SELECT \\* FROM expenses WHERE id=\\$1").WithArgs(1).WillReturnRows(rows)

		if err = GetByIdHandler(c); err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		var expen Expense
		json.NewDecoder(resp.Body).Decode(&expen)
		assert.Equal(t, 1, expen.ID)
		assert.Equal(t, "strawberry smoothie", expen.Title)
		assert.Equal(t, float64(79), expen.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", expen.Note)
		assert.Equal(t, []string{"food", "beverage"}, expen.Tags)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

	t.Run("Get-All-Expense", func(t *testing.T) {
		e := echo.New()
		req, err := http.NewRequest(http.MethodGet, "/expenses", nil)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when making API request", err)
		}
		//req.Header.Set(echo.HeaderAuthorization, "November 10, 2009")
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		c.SetPath("/expenses")

		conn, mock, err := sqlmock.New()
		db = conn
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(1, "apple smoothie", 89, "no discount", pq.Array([]string{"beverage"})).
			AddRow(2, "iPhone 14 Pro Max 1TB", 66900, "birthday gift from my love", pq.Array([]string{"gadget"}))

		mock.ExpectQuery("SELECT \\* FROM expenses").WillReturnRows(rows)

		if err = GetAllHandler(c); err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		var expenses []Expense
		json.NewDecoder(resp.Body).Decode(&expenses)

		assert.Equal(t, 1, expenses[0].ID)
		assert.Equal(t, "apple smoothie", expenses[0].Title)
		assert.Equal(t, float64(89), expenses[0].Amount)
		assert.Equal(t, "no discount", expenses[0].Note)
		assert.Equal(t, []string{"beverage"}, expenses[0].Tags)

		assert.Equal(t, 2, expenses[1].ID)
		assert.Equal(t, "iPhone 14 Pro Max 1TB", expenses[1].Title)
		assert.Equal(t, float64(66900), expenses[1].Amount)
		assert.Equal(t, "birthday gift from my love", expenses[1].Note)
		assert.Equal(t, []string{"gadget"}, expenses[1].Tags)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

}
