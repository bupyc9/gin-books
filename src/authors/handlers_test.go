package authors_test

import (
	"books/authors"
	"books/database"
	"books/router"
	"fmt"

	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
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
	db := database.CreateDb()
	r := router.SetupRouter(db)
	w := httptest.NewRecorder()

	author := authors.Author{
		FirstName:  "First Name",
		LastName:   "Last Name",
		SecondName: "Second Name",
	}
	db.Create(&author)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Nil(t, err)
	assert.ElementsMatch(
		t,
		[]string{
			"id",
			"createdAt",
			"updatedAt",
			"firstName",
			"lastName",
			"secondName",
		},
		slices.Collect(maps.Keys(responseBody)),
	)
	assert.Equal(t, "First Name", responseBody["firstName"])
	assert.Equal(t, "Last Name", responseBody["lastName"])
	assert.Equal(t, "Second Name", responseBody["secondName"])
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
