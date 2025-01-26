package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupRouter(db *sql.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router, db, rdb)
	return router
}

type CreateTweetRequest struct {
	Content string `json:"content"`
	UserID  int64  `json:"authorId"`
}

type CreateTweetResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func TestTweetCreation(t *testing.T) {
	// Mock database and Redis
	db, _, err := getMockDB() // Replace with a mock DB or use sqlmock
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	rdb := getMockRedis() // Replace with a mock Redis

	// Setup router
	router := setupRouter(db, rdb)

	// Primero, crea un usuario con el formato correcto
	userPayload := map[string]interface{}{
		"name":     "Mauricio Giaconia",
		"email":    "maurigiaconia@hotmail.com",
		"password": "1223",
	}

	userRequestBody, _ := json.Marshal(userPayload)
	req, _ := http.NewRequest("POST", "/users/create", bytes.NewBuffer(userRequestBody))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	assert.Equal(t, http.StatusCreated, w.Code)

	requests := []struct {
		payload  CreateTweetRequest
		expected int    // Código de estado esperado
		message  string // Mensaje esperado en la respuesta
	}{
		{CreateTweetRequest{Content: "test posteado", UserID: 1}, int(http.StatusCreated), "Tweet posted"},
		{CreateTweetRequest{Content: "Uala ha transformado la experiencia financiera de millones de usuarios al ofrecer una plataforma accesible y completa que les permite gestionar su dinero, realizar pagos, ahorrar, solicitar créditos y acceder a una amplia gama de servicios financieros con solo unos clics, facilitando su día a día.", UserID: 1}, int(http.StatusBadRequest), "The content of the tweet must not exceed 280 characters"},
		{CreateTweetRequest{Content: "test", UserID: 9999}, int(http.StatusBadRequest), "Nonexistent user"},
	}
	for i, tc := range requests {
		t.Run(fmt.Sprintf("Request %d", i+1), func(t *testing.T) {
			// Crea el cuerpo de la solicitud
			requestBody, _ := json.Marshal(tc.payload)

			fmt.Println(tc.payload)

			// Crea la solicitud HTTP
			req, _ := http.NewRequest("POST", "/tweets/create", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Ejecuta la solicitud
			router.ServeHTTP(w, req)

			// Valida el código de estado
			assert.Equal(t, tc.expected, w.Code)

			// Valida el mensaje de la respuesta si es aplicable
			var response CreateTweetResponse

			if w.Code == http.StatusCreated {
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.message, response.Data)
			} else {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err, "Error deserializando el cuerpo de la respuesta con error")
				assert.Equal(t, tc.expected, errorResponse.Code)
				assert.Contains(t, errorResponse.Error, tc.message)
			}
		})
	}

}

func getMockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	// currentTime := time.Now().Format("2006-01-02 15:04:05")

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("Mauricio Giaconia", "maurigiaconia@hotmail.com", "1223").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectExec("INSERT INTO users").
		WithArgs("Mauricio Giaconia", "maurigiaconia@hotmail.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT id, name, email, created_at FROM users WHERE id = \\$1").
		WithArgs(1). // El ID que estamos buscando
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "created_at"}).
			AddRow(1, "Mauricio Giaconia", "maurigiaconia@hotmail.com", time.Now())) // Datos del usuario con ID 1

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tweets").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	return db, mock, nil
}

func getMockRedis() *redis.Client {
	// Puedes simular el comportamiento de Redis si es necesario
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Dirección de Redis para pruebas
	})
	return client
}
