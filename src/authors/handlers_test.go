package authors_test

import (
	"books/authors"
	"books/router"
	"books/tests"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthorsTestSuite struct {
	tests.TestSuite
}

func (suite *AuthorsTestSuite) TestAuthorCreate() {
	createAuthor := authors.CreateAuthor{
		FirstName:  "First Name",
		LastName:   "Last Name",
		SecondName: "Second Name",
	}
	body, err := json.Marshal(createAuthor)
	req, err := http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusCreated, w.Code)

	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), author.CreatedAt)
	assert.NotEmpty(suite.T(), author.UpdatedAt)
	assert.Equal(suite.T(), "First Name", author.FirstName)
	assert.Equal(suite.T(), "Last Name", author.LastName)
	assert.Equal(suite.T(), "Second Name", author.SecondName)

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(suite.T(), err)
	w = httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *AuthorsTestSuite) TestAuthorCreateValidation() {
	testList := []struct {
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
		{
			"invalid type firstName",
			map[string]any{"firstName": 1},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"firstName": "The expected type is \"string\"",
			}},
		},
		{
			"invalid type lastName",
			map[string]any{"lastName": 1},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"lastName": "The expected type is \"string\"",
			}},
		},
		{
			"invalid type secondName",
			map[string]any{"secondName": 1},
			router.ValidationResponse{Message: "Validation Error", Errors: map[string]string{
				"secondName": "The expected type is \"string\"",
			}},
		},
	}

	for _, tt := range testList {
		suite.T().Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			suite.Router.ServeHTTP(w, req)

			require.Equal(t, http.StatusUnprocessableEntity, w.Code)
			responseJson, _ := json.Marshal(tt.responseBody)
			assert.JSONEq(t, string(responseJson), w.Body.String())
		})
	}
}

func (suite *AuthorsTestSuite) TestAuthorFind() {
	createAuthor := authors.CreateAuthor{
		FirstName:  "First Name",
		LastName:   "Last Name",
		SecondName: "Second Name",
	}
	body, err := json.Marshal(createAuthor)
	req, err := http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusCreated, w.Code)

	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	require.NoError(suite.T(), err)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	w = httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &author)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), author.CreatedAt)
	assert.NotEmpty(suite.T(), author.UpdatedAt)
	assert.Equal(suite.T(), "First Name", author.FirstName)
	assert.Equal(suite.T(), "Last Name", author.LastName)
	assert.Equal(suite.T(), "Second Name", author.SecondName)
}

func (suite *AuthorsTestSuite) TestAuthorFindNotFound() {
	req, _ := http.NewRequest("GET", "/api/authors/100500", nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusNotFound, w.Code)
	responseJson, _ := json.Marshal(router.MessageResponse{Message: "record not found"})
	assert.JSONEq(suite.T(), string(responseJson), w.Body.String())
}

func (suite *AuthorsTestSuite) TestAuthorDelete() {
	createAuthor := authors.CreateAuthor{
		FirstName:  "First Name",
		LastName:   "Last Name",
		SecondName: "Second Name",
	}
	body, err := json.Marshal(createAuthor)
	req, err := http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusCreated, w.Code)

	var author authors.Author
	err = json.Unmarshal(w.Body.Bytes(), &author)
	require.NoError(suite.T(), err)

	req, err = http.NewRequest("DELETE", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(suite.T(), err)
	w = httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusNoContent, w.Code)

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/authors/%d", author.ID), nil)
	assert.NoError(suite.T(), err)
	w = httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *AuthorsTestSuite) TestAuthorDeleteNotFound() {
	req, err := http.NewRequest("DELETE", "/api/authors/100500", nil)
	assert.NoError(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	require.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *AuthorsTestSuite) TestAuthorList() {
	users := []authors.Author{
		{FirstName: "First Name 1", LastName: "Last Name 1"},
		{FirstName: "First Name 2", LastName: "Last Name 2"},
		{FirstName: "First Name 3", LastName: "Last Name 3"},
		{FirstName: "First Name 4", LastName: "Last Name 4"},
		{FirstName: "First Name 5", LastName: "Last Name 5"},
	}
	result := suite.DB.Create(&users)
	require.NoError(suite.T(), result.Error)

	testList := []struct {
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

	for _, tt := range testList {
		suite.T().Run(tt.name, func(t *testing.T) {
			endpoint, _ := url.Parse("/api/authors")
			query := url.Values{}
			for key, value := range tt.requestBody {
				query.Add(key, value.(string))
			}
			endpoint.RawQuery = query.Encode()

			req, _ := http.NewRequest("GET", endpoint.String(), nil)
			w := httptest.NewRecorder()
			suite.Router.ServeHTTP(w, req)

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

func TestSuite(t *testing.T) {
	suite.Run(t, new(AuthorsTestSuite))
}
