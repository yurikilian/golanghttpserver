package exception

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertProblem(t *testing.T, expected Problem, actual Problem) {
	assert.Equal(t, http.StatusBadRequest, actual.Code)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Message, actual.Message)
	assert.Equal(t, expected.Instance, actual.Instance)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.FieldErrors, actual.FieldErrors)
}
