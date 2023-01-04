//go:build integration
// +build integration

package expense

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

const serverPort = 2565

func TestCreateHandler(t *testing.T) {

	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		// h := NewApplication(db)

		e.POST("/expenses", expense.CreateHandler)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// // Arrange
	// reqBody := ``
	// req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/news", serverPort), strings.NewReader(reqBody))
	// assert.NoError(t, err)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// client := http.Client{}

	// // Act
	// resp, err := client.Do(req)
	// assert.NoError(t, err)

	// byteBody, err := ioutil.ReadAll(resp.Body)
	// assert.NoError(t, err)
	// resp.Body.Close()

	// // Assertions
	// expected := "[{\"ID\":1,\"Title\":\"test-title\",\"Content\":\"test-content\",\"Author\":\"test-author\"}]"

	// if assert.NoError(t, err) {
	// 	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// 	assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// err = eh.Shutdown(ctx)
	// assert.NoError(t, err)

	// eh := setupServer()

	// reqBody := bytes.NewBufferString(`{
	// 	"title": "strawberry smoothie",
	// 	"amount": 79,
	// 	"note": "night market promotion discount 10 bath",
	// 	"tags": ["food", "beverage"]
	// }`)
	// var e Expense

	// res := request(http.MethodPost, uri("expenses"), reqBody)
	// err := res.Decode(&e)

	// //assert.Nil(t, err)
	// assert.Equal(t, http.StatusCreated, res.StatusCode)
	// //assert.NotEqual(t, 0, e.ID)
	// assert.Equal(t, "strawberry smoothie", e.Title)
	// //assert.Equal(t, 79.0, e.Amount)
	// //assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	// //assert.Equal(t, []string{"food", "beverage"}, e.Tags)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// err = eh.Shutdown(ctx)
	// assert.NoError(t, err)
}

func uri(paths ...string) string {
	host := "http://localhost" + ":2565"
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
	serverPort := ":2565"
	eh := echo.New()
	go func(e *echo.Echo) {
		InitDB()
		e.POST("/expenses", CreateHandler)
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
