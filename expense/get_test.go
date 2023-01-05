package expense

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetById(t *testing.T) {
	t.Run("Get-Expense-by-ID", func(t *testing.T) {
		expen := Expense{
			ID:     1,
			Title:  "",
			Amount: 0,
			Note:   "",
			Tags:   []string{},
		}
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(1, "strawberry smoothie", 79, "night market promotion discount 10 bath", pq.Array([]string{"food", "beverage"}))
		mock.ExpectQuery("^SELECT \\* FROM expenses WHERE id=\\$1").WithArgs(expen.ID).WillReturnRows(rows)

		if err = getByID(db, &expen); err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}

		assert.Nil(t, err)
		assert.Equal(t, 1, expen.ID)
		assert.Equal(t, "strawberry smoothie", expen.Title)
		assert.Equal(t, 79.0, expen.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", expen.Note)
		assert.Equal(t, []string{"food", "beverage"}, expen.Tags)
	})

}
