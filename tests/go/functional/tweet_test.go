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

type CreateTweetRequest struct {
	Content string `json:"content"`
	UserID  int64  `json:"authorId"`
}

type CreateTweetResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

type ErrorTweetResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func TestTweetCreation(t *testing.T) {
	// Se reutiliza sqlite en memoria para los tests
	db, err := factory.GetDatabase("sqlite")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}
	rdb := getMockRedis()

	conn, err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect DB: %v", err)
	}

	defer conn.Close()

	router := setupTweetRouter(conn, rdb)

	// Se crea un user para realizar el test de tweetear
	userPayload := map[string]interface{}{
		"name":     "Mauricio Giaconia",
		"email":    "maurigiaconia@hotmail.com",
		"password": "1223",
	}

	w := makeRequest(t, "POST", "/users/create", userPayload, router)

	assert.Equal(t, http.StatusCreated, w.Code)

	// []struct con el lsitado de pruebas a realizar
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

			requestBody, _ := json.Marshal(tc.payload)

			req, _ := http.NewRequest("POST", "/tweets/create", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expected, w.Code)

			var response CreateTweetResponse

			if w.Code == http.StatusCreated {
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.message, response.Data)
			} else {
				var errorResponse ErrorTweetResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err, "Error deserializando el cuerpo de la respuesta con error")
				assert.Equal(t, tc.expected, errorResponse.Code)
				assert.Contains(t, errorResponse.Error, tc.message)
			}
		})
	}
}

func TestGetTimeline(t *testing.T) {
	// Se reutiliza sqlite en memoria para los tests
	db, err := factory.GetDatabase("sqlite")
	if err != nil {
		t.Fatalf("failed to create DB: %v", err)
	}
	rdb := getMockRedis()

	conn, err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect DB: %v", err)
	}

	defer conn.Close()

	router := setupTweetRouter(conn, rdb)

	// Se crea un user para realizar el test de tweetear
	userPayload := map[string]interface{}{
		"name":     "Mauricio Giaconia",
		"email":    "maurigiaconia@hotmail.com",
		"password": "1223",
	}

	w := makeRequest(t, "POST", "/users/create", userPayload, router)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Crear el segundo usuario
	otherUserPayload := map[string]interface{}{
		"name":     "Juan Perez",
		"email":    "juanperez@hotmail.com",
		"password": "4567",
	}
	w = makeRequest(t, "POST", "/users/create", otherUserPayload, router)
	assert.Equal(t, http.StatusCreated, w.Code)

	// El usuario 1 sigue al usuario 2
	followPayload := map[string]interface{}{
		"followerId": 2,
		"followedId": 1,
	}
	w = makeRequest(t, "POST", "/users_follow/create", followPayload, router)
	assert.Equal(t, http.StatusOK, w.Code)

	tweetPayload := CreateTweetRequest{Content: "Tweet desde test 1", UserID: 1}
	w = makeRequest(t, "POST", "/tweets/create", tweetPayload, router)
	assert.Equal(t, http.StatusCreated, w.Code)

	tweetPayload = CreateTweetRequest{Content: "Tweet Uala desde test 2", UserID: 1}
	w = makeRequest(t, "POST", "/tweets/create", tweetPayload, router)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = makeRequest(t, "GET", "/tweets/2/timeline", nil, router)
	assert.Equal(t, http.StatusOK, w.Code)

	// Asegurarse de que los tweets de Juan Perez aparecen en la timeline de Mauricio Giaconia
	var response []CreateTweetResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Contains(t, response[0].Data, "Tweet desde test 1")
	assert.Contains(t, response[1].Data, "Tweet Uala desde test 2")
}

func getMockRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return client
}

func setupTweetRouter(db *sql.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router, db, rdb)
	return router
}

// Funcion utilziada para realizar requests necesarias para el test (por ejemplo, si se necesita crear un usuario para poder testear los endpoints de tweets)
func makeRequest(t *testing.T, method, url string, body interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	var requestBody []byte
	if body != nil {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Error marshalling request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	return w
}
