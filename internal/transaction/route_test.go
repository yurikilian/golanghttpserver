package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/yurikilian/bills/internal/logger"
	"github.com/yurikilian/bills/pkg/exception"
	"github.com/yurikilian/bills/pkg/middleware"
	"github.com/yurikilian/bills/pkg/server"
	"github.com/yurikilian/bills/pkg/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Transaction_Create(t *testing.T) {

	tests := []struct {
		name                string
		request             *CreationRequest
		expectedStatusCode  int
		expectsErr          bool
		exceptedEx          exception.Problem
		expectedSavedEntity *Entity
	}{
		{
			name: "Should return 204 no content given correct creation request payload",
			request: &CreationRequest{
				Title:       "Supermarket",
				Description: "Mensal shop",
				Currency:    "EUR",
				Type:        "CREDIT",
				Price:       53.25,
			},
			expectedStatusCode: http.StatusNoContent,
			expectedSavedEntity: &Entity{
				Id:          1.0,
				Title:       "Supermarket",
				Description: "Mensal shop",
				Currency:    "EUR",
				Type:        "CREDIT",
				Price:       53.25,
			},
		},
		{
			name: "Should return 400 bad request given invalid creation request payload with wrong currency",
			request: &CreationRequest{
				Title:       "",
				Description: "",
				Currency:    "DDD",
				Type:        "",
			},
			expectsErr:         true,
			expectedStatusCode: http.StatusBadRequest,
			exceptedEx: exception.NewValidationProblem(
				[]exception.ValidationProblemDetail{
					exception.NewValidationProblemDetail("required", "Title", ""),
					exception.NewValidationProblemDetail("required", "Description", ""),
					exception.NewValidationProblemDetail("required", "Price", ""),
					exception.NewValidationProblemDetail("oneof", "Currency", "EUR"),
					exception.NewValidationProblemDetail("required", "Type", ""),
				},
			),
		},
		{
			name:               "Should return 400 bad request given invalid creation request payload with empty currency",
			request:            &CreationRequest{},
			expectedStatusCode: http.StatusBadRequest,
			expectsErr:         true,

			exceptedEx: exception.NewValidationProblem(
				[]exception.ValidationProblemDetail{
					exception.NewValidationProblemDetail("required", "Title", ""),
					exception.NewValidationProblemDetail("required", "Description", ""),
					exception.NewValidationProblemDetail("required", "Price", ""),
					exception.NewValidationProblemDetail("required", "Currency", ""),
					exception.NewValidationProblemDetail("required", "Type", ""),
				},
			),
		},
		{
			name:               "Should return 400 bad request given invalid creation request payload with nil fields",
			request:            &CreationRequest{},
			expectedStatusCode: http.StatusBadRequest,
			expectsErr:         true,

			exceptedEx: exception.NewValidationProblem(
				[]exception.ValidationProblemDetail{
					exception.NewValidationProblemDetail("required", "Title", ""),
					exception.NewValidationProblemDetail("required", "Description", ""),
					exception.NewValidationProblemDetail("required", "Price", ""),
					exception.NewValidationProblemDetail("required", "Currency", ""),
					exception.NewValidationProblemDetail("required", "Type", ""),
				},
			),
		},
		{
			name: "Should return 400 bad request given invalid transaction type",
			request: &CreationRequest{
				Title:       "Supermarket",
				Description: "Mensal shop",
				Currency:    "EUR",
				Type:        "DREBIT",
				Price:       53.25,
			},
			expectsErr: true,

			expectedStatusCode: http.StatusBadRequest,
			exceptedEx: exception.NewValidationProblem(
				[]exception.ValidationProblemDetail{
					exception.NewValidationProblemDetail("oneof", "Type", "'CREDIT' 'DEBIT'"),
				},
			),
		},
	}

	inMemoryDb := storage.NewInMemoryStorage[Entity]()
	moduleProvider := NewTransactionModuleBuilder().WithInMemoryStorage(inMemoryDb).Build()

	router := server.NewRestRouter().
		Get("/transactions", moduleProvider.ProvideRoute().Find).
		POST("/transactions", moduleProvider.ProvideRoute().Create)

	restServer := server.NewRestServer(server.NewRestServerOptions(":3050", logger.NewProvider().ProvideLog())).
		Router(router).
		Use(middleware.Otel()).
		Use(middleware.Json())

	defer func() {
		_, _ = restServer.Start(nil)
	}()

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			jsonReq, err := json.Marshal(test.request)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(jsonReq))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			restServer.ServeHTTP(rec, req)

			assert.Equal(t, test.expectedStatusCode, rec.Code)

			if test.expectsErr {
				var problem exception.Problem
				err := json.Unmarshal(rec.Body.Bytes(), &problem)
				if err != nil {
					assert.NoError(t, err)
				}

				exception.AssertProblem(t, test.exceptedEx, problem)
			}

			if test.expectedSavedEntity != nil {
				saved, err := inMemoryDb.Find(test.expectedSavedEntity.Id)
				assert.NoError(t, err)
				assert.NotNil(t, saved)
			}

		})
	}

	restServer.Shutdown(context.Background())
}
