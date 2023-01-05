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

func TestUnitCreate(t *testing.T) {
	t.Run("Create expense", func(t *testing.T) {
		e := echo.New()
		reqBody := `{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath",
			"tags": ["food", "beverage"]
		}`
		req, err := http.NewRequest(http.MethodPost, "/expenses", strings.NewReader(reqBody))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when making API request", err)
		}

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

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(1)

		mock.ExpectQuery("INSERT INTO expenses").
			WithArgs("strawberry smoothie",
				float64(79),
				"night market promotion discount 10 bath",
				pq.Array([]string{"food", "beverage"})).
			WillReturnRows(rows)

		if err = CreateHandler(c); err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)

		var expen Expense
		json.NewDecoder(resp.Body).Decode(&expen)

		assert.NotEqual(t, 0, expen.ID)
		assert.Equal(t, "strawberry smoothie", expen.Title)
		assert.Equal(t, float64(79), expen.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", expen.Note)
		assert.Equal(t, []string{"food", "beverage"}, expen.Tags)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expections: %s", err)
		}
	})

}
