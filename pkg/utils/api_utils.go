package utils

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"
)

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"` //Data es de tipo interface, esto se debe a que 'interface' puede contener cualquier tipo de dato y en "data" puede venir un json o un array
}

type SuccessListResponse struct {
	Code     int           `json:"code"`
	Data     []interface{} `json:"data"`
	Count    int           `json:"count"`    // Total de elementos existentes en la BD
	Limit    int           `json:"limit"`    // Maximo de elementos que se obtiene por pagina
	Offset   int           `json:"offset"`   // Pagina actual
	Next     string        `json:"next"`     // Campo para saber cual es el siguiente indce a consultar, sirve para el paginado
	Previous string        `json:"previous"` // Campo para saber el indice anterior a consultar, sirve para paginado
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type ApiCallOptions struct {
	Headers Headers
	Method  string
	Timeout int
	Body    []byte
}

type Headers struct {
	AuthToken   string
	ContentType string
}

func ApiCall(url string, opts ApiCallOptions) (interface{}, error) {

	//Metodo default en caso de que no se especifique
	if opts.Method == "" {
		opts.Method = "GET"
	}

	//Timeout por defecto en caso de que venga con un valor de 0
	if opts.Timeout == 0 {
		opts.Timeout = 5000
	}

	// Crear un cliente HTTP con timeout
	client := &http.Client{
		Timeout: time.Duration(opts.Timeout) * time.Millisecond,
	}

	// Crear la solicitud HTTP
	req, err := http.NewRequest(opts.Method, url, bytes.NewBuffer(opts.Body))
	if err != nil {
		return buildErrorResponse(500, "Failed to create request"), err
	}

	// Agregar los headers
	req.Header.Set("auth_token", "Bearer "+opts.Headers.AuthToken)
	if opts.Headers.ContentType != "" {
		req.Header.Set("Content-Type", opts.Headers.ContentType)
	} else {
		req.Header.Set("Content-Type", "application/json") // Valor por defecto
	}

	resp, err := client.Do(req)
	if err != nil {
		return buildErrorResponse(500, "Failed to execute request"), err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return buildErrorResponse(500, "Failed to read response body"), err
	}

	// Verifica el código de estado de la respuesta
	if resp.StatusCode >= 400 {
		return buildErrorResponse(resp.StatusCode, "API call failed"), errors.New(string(body))
	}

	//Construyo la respueste exitosa para retornar
	successResponse := SuccessResponse{
		Code: 200,
		Data: string(body),
	}

	return successResponse, nil
}

// Metodo encargado de construir la respuesta que tendrán los distintos endpoints.
// Para una respuesta de un unico dato exitoso, bastará con enviar como argumento "code" y "data".
// En cambio, si la respuesta será un listado que requiere paginado, se deberán mandar todos los argumentos permitidos por la funcion
func ResponseToApi(code int, data interface{}, isAList bool, count int, limit int, offset int) interface{} {
	if code >= 400 {
		return buildErrorResponse(code, data)
	}

	if isAList {
		return buildSuccessListResponse(code, data, count, limit, offset)
	}

	successResponse := SuccessResponse{
		Code: code,
		Data: data,
	}

	return successResponse
}

func buildErrorResponse(code int, data interface{}) ErrorResponse {
	if str, ok := data.(string); ok {
		return ErrorResponse{
			Code:  code,
			Error: str,
		}
	}
	return ErrorResponse{
		Code:  code,
		Error: "Sorry, we had trouble processing the information",
	}
}

func buildSuccessListResponse(code int, data interface{}, count int, limit int, offset int) interface{} {

	if slice, ok := data.([]interface{}); ok {
		return SuccessListResponse{
			Code:     code,
			Data:     slice,
			Count:    count,
			Limit:    limit,
			Offset:   offset,
			Next:     "",
			Previous: "",
		}
	}

	return ErrorResponse{
		Code:  code,
		Error: "Sorry, we had trouble processing the information",
	}
}
