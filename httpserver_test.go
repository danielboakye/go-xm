package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danielboakye/go-xm/config"
	"github.com/danielboakye/go-xm/handlers"
	"github.com/danielboakye/go-xm/helpers"
	"github.com/danielboakye/go-xm/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cfg = config.Configurations{
	DBName:              "xmdb",
	DBUser:              "postgres",
	DBPass:              "postgres",
	AccessTokenDuration: 3600000000000,
	JWTSecretKey:        "91ODUHa4O0wRJCbeCXK4igB4ny/JnFH4jJiMmjf5loo=",
	HTTPPort:            "8080",
	KafkaURL:            "kafka:9092",
}

var ErrDBConnection = errors.New("db connection error")

func newTestHTTPHandler(t *testing.T) (*assert.Assertions, *require.Assertions, sqlmock.Sqlmock, IHTTPHandler) {
	assert := assert.New(t)
	require := require.New(t)

	db, mockDB, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(err)

	validator, err := helpers.NewValidation()
	require.NoError(err)

	var (
		testRepo    = repo.NewRepository(db)
		testHandler = handlers.NewHandler(testRepo, validator, cfg)
	)

	testHTTPHandler := newHTTPHandler(testHandler, cfg)

	return assert, require, mockDB, testHTTPHandler
}

func TestCompany(t *testing.T) {

	t.Run("Create Company", func(t *testing.T) {
		t.Run("Invalid Token", func(t *testing.T) {
			assert, _, _, testHTTPHandler := newTestHTTPHandler(t)

			const invalidAccessToken = "invalid token"

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/api/v1/company", strings.NewReader(`
				{
					"name": "example",
					"description": "company description",
					"amountOfEmployees": 2,
					"registered": false,
					"companyType": "Non Profit"
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", invalidAccessToken))

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			if body, err := io.ReadAll(resp.Body); assert.NoError(err) {
				assert.Equal(`{"error":"unauthorized"}`, string(body))
			}

			assert.Equal(401, resp.StatusCode)

		})

		t.Run("Invalid Parameter", func(t *testing.T) {
			assert, require, _, testHTTPHandler := newTestHTTPHandler(t)

			accessToken, err := helpers.GenerateAccessToken(cfg, testUUID)
			require.NoError(err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/api/v1/company", strings.NewReader(`
				{
					"name": "example",
					"description": "company description",
					"amountOfEmployees": 2,
					"registered": false,
					"companyType": "invalid type"
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			if body, err := io.ReadAll(resp.Body); assert.NoError(err) {
				assert.Equal(`{"error":"invalid parameters"}`, string(body))
			}

			assert.Equal(422, resp.StatusCode)

		})

		t.Run("Duplicate company", func(t *testing.T) {
			assert, require, mockDB, testHTTPHandler := newTestHTTPHandler(t)

			mockDB.ExpectQuery(`
				SELECT
					company_id, company_name, description, amount_of_employees, is_registered, company_type
				FROM companies
				WHERE company_name = $1
					AND deleted_at IS NULL
			`).
				WithArgs("example").
				WillReturnRows(
					sqlmock.NewRows(
						[]string{
							"company_id", "company_name",
							"description", "amount_of_employees",
							"is_registered", "company_type",
						},
					).
						FromCSVString("2899bacc-7107-4cd4-9364-6a6fc4fc2fd3,example,company description,2,false,Non Profit"),
				)

			accessToken, err := helpers.GenerateAccessToken(cfg, testUUID)
			require.NoError(err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/api/v1/company", strings.NewReader(`
				{
					"name": "example",
					"description": "company description",
					"amountOfEmployees": 2,
					"registered": false,
					"companyType": "Non Profit"
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			if body, err := io.ReadAll(resp.Body); assert.NoError(err) {
				assert.Equal(`{"error":"duplicate record"}`, string(body))
			}

			assert.Equal(409, resp.StatusCode)

			assert.NoError(mockDB.ExpectationsWereMet())
		})

		t.Run("success", func(t *testing.T) {
			assert, require, mockDB, testHTTPHandler := newTestHTTPHandler(t)

			mockDB.ExpectQuery(`
				SELECT
					company_id, company_name, description, amount_of_employees, is_registered, company_type
				FROM companies
				WHERE company_name = $1
					AND deleted_at IS NULL
			`).
				WithArgs("example12").
				WillReturnError(sql.ErrNoRows)

			mockDB.ExpectQuery(`
				INSERT INTO companies (company_name, description, amount_of_employees, is_registered, company_type)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING company_id
			`).
				WithArgs("example12", "company description", 2, false, "Non Profit").
				WillReturnRows(
					sqlmock.NewRows([]string{"company_id"}).
						AddRow("ae17b2e2-6b87-4c5b-9c94-3623dacf113b"),
				)

			accessToken, err := helpers.GenerateAccessToken(cfg, testUUID)
			require.NoError(err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "http:/api/v1/company", strings.NewReader(`
				{
					"name": "example12",
					"description": "company description",
					"amountOfEmployees": 2,
					"registered": false,
					"companyType": "Non Profit"
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			if body, err := io.ReadAll(resp.Body); assert.NoError(err) {
				log.Println(string(body))
				// assert.Equal(`{"id":"ae17b2e2-6b87-4c5b-9c94-3623dacf113b","name":"example12","description":"company description","amountOfEmployees":2,"registered":false,"companyType":"Non Profit"}`, string(body))
			}
			log.Println(resp.StatusCode)
			assert.Equal(200, resp.StatusCode)

			assert.NoError(mockDB.ExpectationsWereMet())
		})
	})

	t.Run("Update Company", func(t *testing.T) {

		t.Run("Invalid Token", func(t *testing.T) {
			assert, _, _, testHTTPHandler := newTestHTTPHandler(t)

			const invalidAccessToken = "invalid token"

			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "http:/api/v1/company/ae17b2e2-6b87-4c5b-9c94-3623dacf113b", strings.NewReader(`
				{
					"description": "company description",
					"amountOfEmployees": 2
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", invalidAccessToken))

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			if body, err := io.ReadAll(resp.Body); assert.NoError(err) {
				assert.Equal(`{"error":"unauthorized"}`, string(body))
			}

			assert.Equal(401, resp.StatusCode)

		})

		t.Run("Invalid Parameter", func(t *testing.T) {
			assert, require, _, testHTTPHandler := newTestHTTPHandler(t)

			accessToken, err := helpers.GenerateAccessToken(cfg, testUUID)
			require.NoError(err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "http:/api/v1/company/ae17b2e2-6b87-4c5b-9c94-3623dacf113b", strings.NewReader(`
				{
					"amountOfEmployees": -1
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			if body, err := io.ReadAll(resp.Body); assert.NoError(err) {
				assert.Equal(`{"error":"invalid parameters"}`, string(body))
			}

			assert.Equal(422, resp.StatusCode)

		})
		t.Run("success", func(t *testing.T) {
			assert, require, mockDB, testHTTPHandler := newTestHTTPHandler(t)

			mockDB.ExpectQuery(`
				SELECT
					company_id, company_name, description, amount_of_employees, is_registered, company_type
				FROM companies
				WHERE company_id = $1
					AND deleted_at IS NULL
			`).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{
							"company_id", "company_name",
							"description", "amount_of_employees",
							"is_registered", "company_type",
						},
					).
						FromCSVString("ae17b2e2-6b87-4c5b-9c94-3623dacf113b,example,company description,3,false,Non Profit"),
				)

			mockDB.ExpectExec(`
				UPDATE companies 
				SET 
					company_name = coalesce($2, company_name), 
					description = coalesce($3, description), 
					amount_of_employees = coalesce($4, amount_of_employees), 
					is_registered = coalesce($5, is_registered), 
					company_type = coalesce($6, company_type),
					modified_at = now()
				WHERE
					company_id = $1
			`).
				WillReturnResult(sqlmock.NewResult(1, 1))

			accessToken, err := helpers.GenerateAccessToken(cfg, testUUID)
			require.NoError(err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "http:/api/v1/company/ae17b2e2-6b87-4c5b-9c94-3623dacf113b", strings.NewReader(`
				{
					"amountOfEmployees": 2
				}
			`))
			r.Header.Set("Content-Type", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			_, err = io.ReadAll(resp.Body)
			assert.NoError(err)

			assert.Equal(200, resp.StatusCode)

			assert.NoError(mockDB.ExpectationsWereMet())
		})
	})

	t.Run("Delete Company", func(t *testing.T) {

		t.Run("Invalid Token", func(t *testing.T) {
			assert, _, _, testHTTPHandler := newTestHTTPHandler(t)

			const invalidAccessToken = "invalid token"

			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "http:/api/v1/company/ae17b2e2-6b87-4c5b-9c94-3623dacf113b", nil)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", invalidAccessToken))

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			if body, err := io.ReadAll(resp.Body); assert.NoError(err) {
				assert.Equal(`{"error":"unauthorized"}`, string(body))
			}

			assert.Equal(401, resp.StatusCode)

		})

		t.Run("success", func(t *testing.T) {
			assert, require, mockDB, testHTTPHandler := newTestHTTPHandler(t)

			mockDB.ExpectExec(`
				UPDATE companies 
				SET 
					deleted_at = now()
				WHERE
					company_id = $1
			`).
				WithArgs("ae17b2e2-6b87-4c5b-9c94-3623dacf113b").
				WillReturnResult(sqlmock.NewResult(1, 1))

			accessToken, err := helpers.GenerateAccessToken(cfg, testUUID)
			require.NoError(err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "http:/api/v1/company/ae17b2e2-6b87-4c5b-9c94-3623dacf113b", nil)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			_, err = io.ReadAll(resp.Body)
			assert.NoError(err)

			assert.Equal(200, resp.StatusCode)

			assert.NoError(mockDB.ExpectationsWereMet())
		})
	})

	t.Run("Get Company", func(t *testing.T) {

		t.Run("No record", func(t *testing.T) {
			assert, require, mockDB, testHTTPHandler := newTestHTTPHandler(t)

			mockDB.ExpectQuery(`
				SELECT
					company_id, company_name, description, amount_of_employees, is_registered, company_type
				FROM companies
				WHERE company_id = $1
					AND deleted_at IS NULL
			`).
				WithArgs("ae17b2e2-6b87-4c5b-9c94-3623dacf113b").
				WillReturnError(sql.ErrNoRows)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http:/api/v1/company/ae17b2e2-6b87-4c5b-9c94-3623dacf113b", nil)
			r.Header.Set("Content-Type", "application/json")

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			_, err := io.ReadAll(resp.Body)

			require.NoError(err)

			assert.Equal(404, resp.StatusCode)

			assert.NoError(mockDB.ExpectationsWereMet())
		})

		t.Run("success", func(t *testing.T) {
			assert, _, mockDB, testHTTPHandler := newTestHTTPHandler(t)

			mockDB.ExpectQuery(`
				SELECT
					company_id, company_name, description, amount_of_employees, is_registered, company_type
				FROM companies
				WHERE company_id = $1
					AND deleted_at IS NULL
			`).
				WithArgs("ae17b2e2-6b87-4c5b-9c94-3623dacf113b").
				WillReturnRows(
					sqlmock.NewRows(
						[]string{
							"company_id", "company_name",
							"description", "amount_of_employees",
							"is_registered", "company_type",
						},
					).
						FromCSVString("ae17b2e2-6b87-4c5b-9c94-3623dacf113b,example,company description,2,false,Non Profit"),
				)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http:/api/v1/company/ae17b2e2-6b87-4c5b-9c94-3623dacf113b", nil)
			r.Header.Set("Content-Type", "application/json")

			testHTTPHandler.ServeHTTP(w, r)

			resp := w.Result()

			if body, err := io.ReadAll(resp.Body); assert.NoError(err) {
				assert.Equal(`{"id":"ae17b2e2-6b87-4c5b-9c94-3623dacf113b","name":"example","description":"company description","amountOfEmployees":2,"registered":false,"companyType":"Non Profit"}`, string(body))
			}

			assert.Equal(200, resp.StatusCode)

			assert.NoError(mockDB.ExpectationsWereMet())
		})
	})

}
