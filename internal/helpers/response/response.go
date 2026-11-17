package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// ErrorResponse тело ответа API с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// ReadBody читает и закрывает тело HTTP-ответа
func ReadBody(t *testing.T, rawResponse *http.Response) []byte {
	t.Helper()

	if rawResponse.Body == nil {
		return nil
	}
	defer rawResponse.Body.Close()

	body, err := io.ReadAll(rawResponse.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")

	return body
}

// ReadBodyWithStatus проверяет HTTP status code и возвращает тело ответа;
// при несовпадении кода тело ответа попадает в сообщение об ошибке
func ReadBodyWithStatus(t *testing.T, rawResponse *http.Response, expectedStatusCode int) []byte {
	t.Helper()

	body := ReadBody(t, rawResponse)
	require.Equalf(t, expectedStatusCode, rawResponse.StatusCode,
		"HTTP status code должен быть %d, получен %d. Тело ответа: %s",
		expectedStatusCode, rawResponse.StatusCode, string(body))

	return body
}

// Decode разбирает JSON-тело ответа в структуру указанного типа
func Decode[T any](t *testing.T, body []byte) T {
	t.Helper()

	var result T
	require.NoErrorf(t, json.Unmarshal(body, &result), "Не удалось разобрать JSON ответа: %s", string(body))

	return result
}

// DecodeError разбирает тело ответа с ошибкой и проверяет, что в нём заполнены error и message
func DecodeError(t *testing.T, body []byte) ErrorResponse {
	t.Helper()

	errorResponse := Decode[ErrorResponse](t, body)
	require.NotEmptyf(t, errorResponse.Error, "В ответе с ошибкой пустое поле error. Тело ответа: %s", string(body))
	require.NotEmptyf(t, errorResponse.Message, "В ответе с ошибкой пустое поле message. Тело ответа: %s", string(body))

	return errorResponse
}
