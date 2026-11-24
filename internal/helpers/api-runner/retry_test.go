package api_runner

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const requestBody = `{"advertisement_id":"1"}`

// newTestClient создаёт клиент с повтором без реальных пауз
func newTestClient() *http.Client {
	return &http.Client{Transport: retryTransport{next: http.DefaultTransport, delay: time.Millisecond}}
}

// Проверяем, что после двух ответов 503 запрос повторяется и тело отправляется заново целиком
func TestRetryOnServiceUnavailable(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, requestBody, string(body), "Тело повторного запроса должно совпадать с исходным")
		if calls.Add(1) < retryAttempts {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	response, err := newTestClient().Post(server.URL, "application/json", strings.NewReader(requestBody))
	require.NoError(t, err)
	defer func() { _ = response.Body.Close() }()

	require.Equal(t, http.StatusCreated, response.StatusCode, "После повторов должен прийти ответ сервера")
	require.EqualValues(t, retryAttempts, calls.Load(), "Количество обращений к серверу")
}

// Проверяем, что постоянный 503 возвращается тесту после последней попытки, а ответы 4xx не повторяются
func TestRetryStopsAfterLastAttemptAndSkipsOtherStatuses(t *testing.T) {
	testCases := []struct {
		name          string
		status        int
		expectedCalls int32
	}{
		{name: "Постоянный 503", status: http.StatusServiceUnavailable, expectedCalls: retryAttempts},
		{name: "Ответ 401 не повторяется", status: http.StatusUnauthorized, expectedCalls: 1},
		{name: "Ответ 404 не повторяется", status: http.StatusNotFound, expectedCalls: 1},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				w.WriteHeader(testCase.status)
			}))
			defer server.Close()

			response, err := newTestClient().Get(server.URL)
			require.NoError(t, err)
			defer func() { _ = response.Body.Close() }()

			require.Equal(t, testCase.status, response.StatusCode, "Статус ответа")
			require.Equal(t, testCase.expectedCalls, calls.Load(), "Количество обращений к серверу")
		})
	}
}
