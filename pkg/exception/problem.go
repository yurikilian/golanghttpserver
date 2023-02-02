package exception

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

const baseUrl = "https://mybils.io"

type Problem struct {
	Code        int      `json:"code"`
	Title       string   `json:"title"`
	Message     string   `json:"detail"`
	Instance    string   `json:"instance"`
	Type        string   `json:"type"`
	FieldErrors []string `json:"fieldErrors,omitempty"`
}

func (p *Problem) Error() string {
	return p.Message
}

var _ error = (*Problem)(nil)

func NewInternalServerError(messages ...string) *Problem {
	var message string

	if len(messages) > 0 {
		message = strings.Join(messages, ".")
	} else {
		message = "An undetermined error was triggered. Please, contact the support team"
	}

	return &Problem{
		Code:     http.StatusInternalServerError,
		Message:  message,
		Instance: "N/A",
		Title:    "Internal server error",
		Type:     fmt.Sprintf("%v/problems/internal-server-error", baseUrl),
	}
}

func NewMalformedRequestProblem() *Problem {
	return &Problem{
		Code:     http.StatusUnprocessableEntity,
		Title:    "Malformed request",
		Message:  "The request is malformed, verify the input or parameters sent",
		Instance: "N/A",
		Type:     fmt.Sprintf("%v/problems/malformed-request", baseUrl),
	}
}

func NewBadRequestProblem(message string) *Problem {
	return &Problem{
		Code:     http.StatusBadRequest,
		Title:    "Invalid request",
		Message:  message,
		Instance: "N/A",
		Type:     fmt.Sprintf("%v/problems/bad-request", baseUrl),
	}
}

func NewUnsupportedMediaType(message string) *Problem {
	return &Problem{
		Code:     http.StatusUnsupportedMediaType,
		Title:    "Invalid request",
		Message:  message,
		Instance: "N/A",
		Type:     fmt.Sprintf("%v/problems/unsupported-media-type", baseUrl),
	}
}

func NewValidationProblem(vErrors []*ValidationProblemDetail) *Problem {
	return &Problem{
		Code:        http.StatusBadRequest,
		Title:       "Invalid request",
		Message:     "The request does not satisfy the validation rules",
		Instance:    "N/A",
		Type:        fmt.Sprintf("%v/problems/invalid-request", baseUrl),
		FieldErrors: mapValidationErrors(vErrors),
	}
}
func NewRouteNotFound(path string) *Problem {
	return &Problem{
		Code:     http.StatusNotFound,
		Title:    "Route not found",
		Message:  fmt.Sprintf("The route `%v` does not exist", path),
		Instance: "N/A",
		Type:     fmt.Sprintf("%v/problems/not-found", baseUrl),
	}
}

func NewMethodNotAllowed(path string, method string) *Problem {
	return &Problem{
		Code:     http.StatusMethodNotAllowed,
		Title:    "Method not allowed",
		Message:  fmt.Sprintf("The method %v is not allowed for route `%v`", method, path),
		Instance: "N/A",
		Type:     fmt.Sprintf("%v/problems/method-not-allowed", baseUrl),
	}
}

func mapValidationErrors(vErrors []*ValidationProblemDetail) []string {

	sErrors := make([]string, 0)

	for _, err := range vErrors {
		switch err.Tag {
		case "required":
			sErrors = append(sErrors, fmt.Sprintf("%s is required", err.Field))
		case "email":
			sErrors = append(sErrors, fmt.Sprintf("%s is not valid email", err.Field))
		case "gte":
			sErrors = append(sErrors, fmt.Sprintf("%s value must be greater than %s", err.Field, err.Param))
		case "lte":
			sErrors = append(sErrors, fmt.Sprintf("%s value must be lower than %s", err.Field, err.Param))
		case "oneof":
			sErrors = append(sErrors, fmt.Sprintf("%s value must be one of the following: %s", err.Field, formatParam(err.Param)))
		}

	}

	return sErrors
}

func formatParam(param string) string {

	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}
	a := strings.FieldsFunc(param, f)
	joined := strings.Join(a, ", ")
	joined = strings.ReplaceAll(joined, "'", "")

	lastI := strings.LastIndex(joined, ",")
	if lastI != -1 {
		return joined[:lastI] + string(" or") + joined[lastI+1:]
	} else {
		return joined
	}

}

type ValidationProblemDetail struct {
	Tag   string
	Field string
	Param string
}

func NewValidationProblemDetail(tag string, field string, param string) *ValidationProblemDetail {
	return &ValidationProblemDetail{
		Tag:   tag,
		Field: field,
		Param: param,
	}
}
