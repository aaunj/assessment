//go:build integration
// +build integration

package expense

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var serverPort = os.Getenv("PORT")

func TestIntegrationExpense(t *testing.T) {
	eh := setupServer()

	t.Run("Create", func(t *testing.T) {
		reqBody := `{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath",
			"tags": ["food","beverage"]
		}`

		resp := request(http.MethodPost, uri("expenses"), strings.NewReader(reqBody))

		var e Expense
		err := resp.Decode(&e)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEqual(t, 0, e.ID)
		assert.Equal(t, "strawberry smoothie", e.Title)
		assert.Equal(t, 79.0, e.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", e.Note)
		assert.Equal(t, []string{"food", "beverage"}, e.Tags)
	})

	t.Run("Get-By-Id", func(t *testing.T) {
		reqBody := ``

		resp := request(http.MethodGet, uri("expenses/1"), strings.NewReader(reqBody))

		var e Expense
		err := resp.Decode(&e)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 1, e.ID)
		assert.Equal(t, "strawberry smoothie", e.Title)
		assert.Equal(t, float64(79), e.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", e.Note)
		assert.Equal(t, []string{"food", "beverage"}, e.Tags)
	})

	t.Run("Update-By-Id", func(t *testing.T) {
		reqBody := `{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`

		resp := request(http.MethodPut, uri("expenses/1"), strings.NewReader(reqBody))

		var e Expense
		err := resp.Decode(&e)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 1, e.ID)
		assert.Equal(t, "apple smoothie", e.Title)
		assert.Equal(t, float64(89), e.Amount)
		assert.Equal(t, "no discount", e.Note)
		assert.Equal(t, []string{"beverage"}, e.Tags)
	})

	t.Run("Get-All", func(t *testing.T) {
		reqBody := `{
			"title": "iPhone 14 Pro Max 1TB",
			"amount": 66900,
			"note": "birthday gift from my love",
			"tags": ["gadget"]
		}`

		request(http.MethodPost, uri("expenses"), strings.NewReader(reqBody))

		reqBody = ``

		resp := request(http.MethodGet, uri("expenses"), strings.NewReader(reqBody))

		var expenses []Expense
		err := resp.Decode(&expenses)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

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
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func uri(paths ...string) string {
	host := "http://localhost" + serverPort
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	//assert.NoError(t, err)
	req.Header.Add("Authorization", "November 10, 2009")
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	//assert.NoError(t, err)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func setupServer() *echo.Echo {
	eh := echo.New()
	go func(e *echo.Echo) {
		InitDB()
		e.POST("/expenses", CreateHandler)
		e.GET("/expenses/:id", GetByIdHandler)
		e.PUT("/expenses/:id", UpdateByIdHandler)
		e.GET("/expenses", GetAllHandler)
		e.Start(serverPort)
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost%s", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	return eh
}
