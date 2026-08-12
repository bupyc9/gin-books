package books_test

import (
	"books/authors"
	"books/books"
	"books/router"
	"books/tests"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BooksTestSuite struct {
	tests.TestSuite
}

func (suite *BooksTestSuite) TestBookCreate() {
	author := authors.Author{FirstName: "FirstName", LastName: "LastName"}
	result := suite.DB.Create(&author)
	require.NoError(suite.T(), result.Error)

	createBooks := books.CreateBookRequest{
		AuthorID: author.ID,
		Name:     "Test Book",
		Pages:    300,
		Year:     2010,
	}

	body, err := json.Marshal(createBooks)
	var req *http.Request
	req, err = http.NewRequest("POST", "/api/books", strings.NewReader(string(body)))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusCreated, w.Code)

	var book books.Book
	err = json.Unmarshal(w.Body.Bytes(), &book)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), book.CreatedAt)
	assert.NotEmpty(suite.T(), book.UpdatedAt)
	assert.Equal(suite.T(), "Test Book", book.Name)
	assert.Equal(suite.T(), uint(300), book.Pages)
	assert.Equal(suite.T(), uint(2010), book.Year)
	assert.Equal(suite.T(), "FirstName", book.Author.FirstName)
	assert.Equal(suite.T(), "LastName", book.Author.LastName)
}

func (suite *BooksTestSuite) TestBookCreateValidation() {
	testList := []struct {
		name         string
		requestBody  map[string]any
		responseBody router.ValidationResponse
	}{
		{
			"empty",
			map[string]any{},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"AuthorID": "required",
				"Name":     "required",
				"Pages":    "required",
				"Year":     "required",
			}},
		},
		{
			"invalid type authorId",
			map[string]any{"authorId": "1"},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"authorId": "The expected type is \"uint\"",
			}},
		},
		{
			"invalid type name",
			map[string]any{"name": 1},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"name": "The expected type is \"string\"",
			}},
		},
		{
			"invalid type pages",
			map[string]any{"pages": "1"},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"pages": "The expected type is \"uint\"",
			}},
		},
		{
			"invalid type year",
			map[string]any{"year": "1"},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"year": "The expected type is \"uint\"",
			}},
		},
		{
			"author not found",
			map[string]any{"authorId": 100500, "name": "Test Book", "pages": 300, "year": 2010},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"AuthorID": "not_found",
			}},
		},
	}

	for _, tt := range testList {
		suite.T().Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/books", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			suite.Router.ServeHTTP(w, req)

			require.Equal(t, http.StatusUnprocessableEntity, w.Code)

			responseJson, _ := json.Marshal(tt.responseBody)
			assert.JSONEq(t, string(responseJson), w.Body.String())
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BooksTestSuite))
}
