package books_test

import (
	"books/authors"
	"books/books"
	"books/router"
	"books/tests"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BooksTestSuite struct {
	tests.TestSuite
}

func (suite *BooksTestSuite) TestBookCreate() {
	author := authors.Author{FirstName: "FirstName", LastName: "LastName"}
	result := suite.DB.Create(&author)
	suite.Require().NoError(result.Error)

	createBooks := books.CreateBookRequest{
		AuthorID: author.ID,
		Name:     "Test Book",
		Pages:    300,
		Year:     2010,
	}

	body, err := json.Marshal(createBooks)
	var req *http.Request
	req, err = http.NewRequest("POST", "/api/books", strings.NewReader(string(body)))
	suite.Assert().NoError(err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusCreated, w.Code)

	var book books.Book
	err = json.Unmarshal(w.Body.Bytes(), &book)
	suite.Require().NoError(err)
	suite.Assert().NotEmpty(book.CreatedAt)
	suite.Assert().NotEmpty(book.UpdatedAt)
	suite.Assert().Equal("Test Book", book.Name)
	suite.Assert().Equal(uint(300), book.Pages)
	suite.Assert().Equal(uint(2010), book.Year)
	suite.Assert().Equal("FirstName", book.Author.FirstName)
	suite.Assert().Equal("LastName", book.Author.LastName)
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

			suite.Require().Equal(http.StatusUnprocessableEntity, w.Code)

			responseJson, _ := json.Marshal(tt.responseBody)
			suite.Assert().JSONEq(string(responseJson), w.Body.String())
		})
	}
}

func (suite *BooksTestSuite) TestBookFind() {
	book := books.Book{
		Name:   "Test Book",
		Pages:  300,
		Year:   2010,
		Author: authors.Author{FirstName: "FirstName", LastName: "LastName"},
	}
	result := suite.DB.Create(&book)
	suite.Require().NoError(result.Error)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/books/%d", book.ID), nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusOK, w.Code)

	var responseBook books.Book
	err := json.Unmarshal(w.Body.Bytes(), &responseBook)
	suite.Require().NoError(err)
	suite.Assert().NotEmpty(responseBook.CreatedAt)
	suite.Assert().NotEmpty(responseBook.UpdatedAt)
	suite.Assert().Equal("Test Book", responseBook.Name)
	suite.Assert().Equal(uint(300), responseBook.Pages)
	suite.Assert().Equal(uint(2010), responseBook.Year)
	suite.Assert().Equal("FirstName", responseBook.Author.FirstName)
	suite.Assert().Equal("LastName", responseBook.Author.LastName)
	suite.Assert().Equal("", responseBook.Author.SecondName)
}

func (suite *BooksTestSuite) TestBookFindNotFound() {
	req, _ := http.NewRequest("GET", "/api/books/100500", nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusNotFound, w.Code)
	responseJson, _ := json.Marshal(router.MessageResponse{Message: "record not found"})
	suite.Assert().JSONEq(string(responseJson), w.Body.String())
}

func (suite *BooksTestSuite) TestBookDelete() {
	book := books.Book{
		Name:   "Test Book",
		Pages:  300,
		Year:   2010,
		Author: authors.Author{FirstName: "FirstName", LastName: "LastName"},
	}
	result := suite.DB.Create(&book)
	suite.Require().NoError(result.Error)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/books/%d", book.ID), nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusOK, w.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/books/%d", book.ID), nil)
	w = httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusNotFound, w.Code)
}

func (suite *BooksTestSuite) TestBookDeleteNotFound() {
	req, _ := http.NewRequest("DELETE", "/api/books/100500", nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	suite.Require().Equal(http.StatusNotFound, w.Code)
	responseJson, _ := json.Marshal(router.MessageResponse{Message: "record not found"})
	suite.Assert().JSONEq(string(responseJson), w.Body.String())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BooksTestSuite))
}
