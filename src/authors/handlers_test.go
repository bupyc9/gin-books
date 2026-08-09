package authors_test

import (
	"books/authors"
	"books/database"
	"books/router"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := database.CreateDb()
	transaction := db.Begin()
	defer transaction.Rollback()

	r := router.SetupRouter(transaction)

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

	require.Equal(t, http.StatusCreated, w.Code)

	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	require.NoError(t, err)
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

	db := database.CreateDb()
	transaction := db.Begin()
	defer transaction.Rollback()

	r := router.SetupRouter(transaction)

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

			require.Equal(t, http.StatusUnprocessableEntity, w.Code)
			responseJson, _ := json.Marshal(tt.responseBody)
			assert.JSONEq(t, string(responseJson), w.Body.String())
		})
	}
}

func TestAuthorFind(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := database.CreateDb()
	transaction := db.Begin()
	defer transaction.Rollback()

	r := router.SetupRouter(transaction)

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

	require.Equal(t, http.StatusCreated, w.Code)
	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	require.NoError(t, err)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &author)
	require.NoError(t, err)
	assert.NotEmpty(t, author.CreatedAt)
	assert.NotEmpty(t, author.UpdatedAt)
	assert.False(t, author.DeletedAt.Valid)
	assert.Equal(t, "First Name", author.FirstName)
	assert.Equal(t, "Last Name", author.LastName)
	assert.Equal(t, "Second Name", author.SecondName)
}

func TestAuthorFindNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := router.SetupRouter(database.CreateDb())

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/authors/100500", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
	responseJson, _ := json.Marshal(router.MessageResponse{Message: "record not found"})
	assert.JSONEq(t, string(responseJson), w.Body.String())
}

func TestAuthorDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := database.CreateDb()
	transaction := db.Begin()
	defer transaction.Rollback()

	r := router.SetupRouter(transaction)

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

	require.Equal(t, http.StatusCreated, w.Code)

	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	require.NoError(t, err)

	req, err = http.NewRequest("DELETE", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestAuthorDeleteNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := router.SetupRouter(database.CreateDb())

	w := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/authors/100500", nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestAuthorList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := database.CreateDb()
	transaction := db.Begin()
	defer transaction.Rollback()

	r := router.SetupRouter(transaction)

	w := httptest.NewRecorder()

	users := []authors.Author{
		{FirstName: "First Name 1", LastName: "Last Name 1"},
		{FirstName: "First Name 2", LastName: "Last Name 2"},
		{FirstName: "First Name 3", LastName: "Last Name 3"},
		{FirstName: "First Name 4", LastName: "Last Name 4"},
		{FirstName: "First Name 5", LastName: "Last Name 5"},
	}
	result := transaction.Create(&users)
	require.NoError(t, result.Error)

	tests := []struct {
		name            string
		requestBody     map[string]any
		expectedAuthors []authors.Author
	}{
		{
			"empty",
			map[string]any{},
			[]authors.Author{
				{FirstName: "First Name 1", LastName: "Last Name 1"},
				{FirstName: "First Name 2", LastName: "Last Name 2"},
				{FirstName: "First Name 3", LastName: "Last Name 3"},
				{FirstName: "First Name 4", LastName: "Last Name 4"},
				{FirstName: "First Name 5", LastName: "Last Name 5"},
			},
		},
		{
			"count 1",
			map[string]any{"count": "1"},
			[]authors.Author{
				{FirstName: "First Name 1", LastName: "Last Name 1"},
			},
		},
		{
			"page 2 count 2",
			map[string]any{"page": "2", "count": "2"},
			[]authors.Author{
				{FirstName: "First Name 3", LastName: "Last Name 3"},
				{FirstName: "First Name 4", LastName: "Last Name 4"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint, _ := url.Parse("/api/authors")
			query := url.Values{}
			for key, value := range tt.requestBody {
				query.Add(key, value.(string))
			}
			endpoint.RawQuery = query.Encode()

			req, _ := http.NewRequest("GET", endpoint.String(), nil)
			w = httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)

			var responseAuthors []authors.Author
			err := json.Unmarshal(w.Body.Bytes(), &responseAuthors)
			require.NoError(t, err)
			require.Equal(t, len(tt.expectedAuthors), len(responseAuthors))
			for i, author := range tt.expectedAuthors {
				assert.Equal(t, author.FirstName, responseAuthors[i].FirstName, fmt.Sprintf("Author #%d", i))
				assert.Equal(t, author.LastName, responseAuthors[i].LastName, fmt.Sprintf("Author #%d", i))
			}
		})
	}
}
