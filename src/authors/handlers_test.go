package authors_test

import (
	"books/authors"
	"books/database"
	"books/router"

	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorCreate(t *testing.T) {
	r := router.SetupRouter(database.CreateDb())

	w := httptest.NewRecorder()

	createAuthor := authors.CreateAuthor{
		FirstName:  "First Name",
		LastName:   "Last Name",
		SecondName: "Second Name",
	}
	body, _ := json.Marshal(createAuthor)
	req, _ := http.NewRequest("POST", "/api/authors", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

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

func TestAuthorCreateValidation(t *testing.T) {
	r := router.SetupRouter(database.CreateDb())

	w := httptest.NewRecorder()

	tests := []struct {
		name         string
		requestBody  map[string]any
		responseBody router.ValidationError
	}{
		{
			"empty",
			map[string]any{},
			router.ValidationError{Message: "Validation Error", Errors: map[string]string{
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
