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
	"github.com/MauricioGiaconia/uala_backend_challenge/pkg/utils"
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

type TimelineTweet struct {
	TweetID    int    `json:"tweetId"`
	AuthorID   int    `json:"authorId"`
	AuthorName string `json:"authorName"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdAt"` // O usa `time.Time` si quieres manejar fechas como objetos
}

type TimelineResponse struct {
	Code   int             `json:"code"`
	Data   []TimelineTweet `json:"data"`
	Count  int             `json:"count"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Next   string          `json:"next"`
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

			assert.Equal(t, int64(tc.expected), int64(w.Code))

			var response CreateTweetResponse

			if w.Code == http.StatusCreated {
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.message, response.Data)
			} else {
				var errorResponse utils.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err, "Error deserializando el cuerpo de la respuesta con error")
				assert.Equal(t, int64(tc.expected), int64(errorResponse.Code))
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
	assert.Equal(t, int64(http.StatusCreated), int64(w.Code))
	assert.Equal(t, int64(http.StatusCreated), int64(w.Code))

	// Crear el segundo usuario
	otherUserPayload := map[string]interface{}{
		"name":     "Juan Perez",
		"email":    "juanperez@hotmail.com",
		"password": "4567",
	}
	w = makeRequest(t, "POST", "/users/create", otherUserPayload, router)
	assert.Equal(t, int64(http.StatusCreated), int64(w.Code))

	// El usuario 1 sigue al usuario 2
	followPayload := map[string]interface{}{
		"followerId": 2,
		"followedId": 1,
	}
	w = makeRequest(t, "POST", "/users_follow/create", followPayload, router)
	assert.Equal(t, int64(http.StatusCreated), int64(w.Code))

	tweetPayload := CreateTweetRequest{Content: "Tweet desde test 1", UserID: 1}
	w = makeRequest(t, "POST", "/tweets/create", tweetPayload, router)
	assert.Equal(t, int64(http.StatusCreated), int64(w.Code))

	tweetPayload = CreateTweetRequest{Content: "Tweet Uala desde test 2", UserID: 1}
	w = makeRequest(t, "POST", "/tweets/create", tweetPayload, router)
	assert.Equal(t, int64(http.StatusCreated), int64(w.Code))

	w = makeRequest(t, "GET", "/tweets/2/timeline", nil, router)
	assert.Equal(t, int64(http.StatusOK), int64(w.Code))

	// Asegurarse de que los tweets de Mauricio Giaconia aparecen en la timeline de Juan Perez
	var response TimelineResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Len(t, response.Data, 2)
	assert.Contains(t, response.Data[0].Content, "Tweet desde test 1")
	assert.Contains(t, response.Data[1].Content, "Tweet Uala desde test 2")
	assert.Equal(t, response.Count, 2)
	assert.Equal(t, response.Limit, 25)
	assert.Equal(t, response.Offset, 0)

	var errorResponse utils.ErrorResponse
	// Obtencion del timeline de un usuario inexistente
	w = makeRequest(t, "GET", "/tweets/999/timeline", nil, router)
	assert.Equal(t, int64(http.StatusBadRequest), int64(w.Code))

	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)

	assert.NoError(t, err)
	assert.Contains(t, errorResponse.Error, "Nonexistent user")
}

func getMockRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
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
