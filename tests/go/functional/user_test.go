package functional

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MauricioGiaconia/uala_backend_challenge/internal/routes"
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/factory"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSuccessResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

type ErrorUserResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func TestUserCreation(t *testing.T) {
	// Se reutiliza sqlite en memoria para los tests
	db, err := factory.GetDatabase("sqlite")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}

	conn, err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect DB: %v", err)
	}

	defer conn.Close()

	router := setupUserRouter(conn, &redis.Client{})

	// Se crea un user para realizar el test de tweetear
	userPayload := map[string]interface{}{
		"name":     "Mauricio Giaconia",
		"email":    "maurigiaconia@hotmail.com",
		"password": "1223",
	}

	userRequestBody, _ := json.Marshal(userPayload)
	req, _ := http.NewRequest("POST", "/users/create", bytes.NewBuffer(userRequestBody))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// []struct con el lsitado de pruebas a realizar
	requests := []struct {
		payload  CreateUserRequest
		expected int    // Código de estado esperado
		data     int    // Mensaje esperado en la respuesta
		errors   string //Mensaje esperado en caso de error
	}{
		{CreateUserRequest{Name: "Mauricio Giaconia", Email: "maurigiaconia@hotmail.com", Password: "1234"}, int(http.StatusCreated), 1, ""},
		{CreateUserRequest{Name: "User Uala", Email: "user@uala.com.ar", Password: "1234"}, int(http.StatusCreated), 1, ""},
		{CreateUserRequest{Name: "Uala Test", Email: "uala@test.com.ar", Password: "1234"}, int(http.StatusCreated), 1, ""},
		{CreateUserRequest{Name: "Uala Test", Email: "user@uala.com.a", Password: "1234"}, int(http.StatusInternalServerError), 0, "Error to create user"},
	}

	for i, tc := range requests {
		t.Run(fmt.Sprintf("Request %d", i+1), func(t *testing.T) {

			requestBody, _ := json.Marshal(tc.payload)

			req, _ := http.NewRequest("POST", "/users/create", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)

			var response UserSuccessResponse

			if w.Code == http.StatusCreated {
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.data, response.Data)
			} else {
				var errorResponse ErrorUserResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err, "Error deserializando el cuerpo de la respuesta con error")
				assert.Equal(t, tc.expected, errorResponse.Code)
				assert.Contains(t, errorResponse.Error, tc.errors)
			}
		})
	}
}

func setupUserRouter(db *sql.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router, db, rdb)
	return router
}
