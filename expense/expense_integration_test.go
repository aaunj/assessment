//go:build integration

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

func TestCreateHandler(t *testing.T) {

	eh := setupServer()

	reqBody := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`
	var e Expense

	res := request(t, http.MethodPost, uri("expenses"), strings.NewReader(reqBody))
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, 79.0, e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func uri(paths ...string) string {
	host := "http://localhost" + os.Getenv("PORT")
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(t *testing.T, method, url string, body io.Reader) *Response {
	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}
	res, err := client.Do(req)
	assert.NoError(t, err)
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
	serverPort := os.Getenv("PORT")
	eh := echo.New()
	go func(e *echo.Echo) {
		InitDB()
		e.GET("/expenses", CreateHandler)
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
