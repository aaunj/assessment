package expense

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("Create expense", func(t *testing.T) {
		expen := Expense{
			ID:     0,
			Title:  "strawberry smoothie",
			Amount: 79,
			Note:   "night market promotion discount 10 bath",
			Tags:   []string{"food", "beverage"},
		}
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("INSERT INTO expenses").WithArgs(
			expen.Title,
			expen.Amount,
			expen.Note,
			pq.Array(expen.Tags)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

		if err = queryRow(db, &expen); err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}

		assert.Nil(t, err)
		assert.NotEqual(t, 0, expen.ID)
		assert.Equal(t, "strawberry smoothie", expen.Title)
		assert.Equal(t, 79.0, expen.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", expen.Note)
		assert.Equal(t, []string{"food", "beverage"}, expen.Tags)
	})

}
