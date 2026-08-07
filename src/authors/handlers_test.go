package authors_test

import (
	"books/authors"
	"books/database"
	"books/router"
	"fmt"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthorCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter(database.CreateDb())

	var w *httptest.ResponseRecorder
	w = httptest.NewRecorder()

	createAuthor := authors.CreateAuthor{
		FirstName:  "First Name",
		LastName:   "Last Name",
		SecondName: "Second Name",
	}
	var err error
	body, err := json.Marshal(createAuthor)
	var req *http.Request
	req, err = http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	assert.NoError(t, err)
	assert.NotEmpty(t, author.CreatedAt)
	assert.NotEmpty(t, author.UpdatedAt)
	assert.False(t, author.DeletedAt.Valid)
	assert.Equal(t, "First Name", author.FirstName)
	assert.Equal(t, "Last Name", author.LastName)
	assert.Equal(t, "Second Name", author.SecondName)

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthorCreateValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter(database.CreateDb())

	w := httptest.NewRecorder()

	tests := []struct {
		name         string
		requestBody  map[string]any
		responseBody router.ValidationResponse
	}{
		{
			"empty",
			map[string]any{},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"FirstName": "required",
				"LastName":  "required",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
			responseJson, _ := json.Marshal(tt.responseBody)
			assert.JSONEq(t, string(responseJson), w.Body.String())
		})
	}
}

func TestAuthorFind(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter(database.CreateDb())
	var w *httptest.ResponseRecorder
	w = httptest.NewRecorder()

	createAuthor := authors.CreateAuthor{
		FirstName:  "First Name",
		LastName:   "Last Name",
		SecondName: "Second Name",
	}
	var err error
	body, err := json.Marshal(createAuthor)
	var req *http.Request
	req, err = http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	assert.NoError(t, err)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &author)
	assert.NoError(t, err)
	assert.NotEmpty(t, author.CreatedAt)
	assert.NotEmpty(t, author.UpdatedAt)
	assert.False(t, author.DeletedAt.Valid)
	assert.Equal(t, "First Name", author.FirstName)
	assert.Equal(t, "Last Name", author.LastName)
	assert.Equal(t, "Second Name", author.SecondName)
}

func TestAuthorFindNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := database.CreateDb()
	r := router.SetupRouter(db)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/authors/100500", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	responseJson, _ := json.Marshal(router.MessageResponse{Message: "record not found"})
	assert.JSONEq(t, string(responseJson), w.Body.String())
}

func TestAuthorDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter(database.CreateDb())

	var w *httptest.ResponseRecorder
	w = httptest.NewRecorder()

	createAuthor := authors.CreateAuthor{
		FirstName:  "First Name",
		LastName:   "Last Name",
		SecondName: "Second Name",
	}
	var err error
	body, err := json.Marshal(createAuthor)
	var req *http.Request
	req, err = http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	assert.NoError(t, err)

	req, err = http.NewRequest("DELETE", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAuthorDeleteNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter(database.CreateDb())

	var w *httptest.ResponseRecorder
	w = httptest.NewRecorder()

	createAuthor := authors.CreateAuthor{
		FirstName:  "First Name",
		LastName:   "Last Name",
		SecondName: "Second Name",
	}
	var err error
	body, err := json.Marshal(createAuthor)
	var req *http.Request
	req, err = http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	assert.NoError(t, err)

	req, err = http.NewRequest("DELETE", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
