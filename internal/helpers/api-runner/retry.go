package api_runner

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"api-tests-template/internal/utils"
)

const (
	// retryAttempts общее число попыток отправить запрос, если стенд отвечает 503
	retryAttempts = 3
	// retryDelay пауза перед второй попыткой, перед третьей она вдвое больше
	retryDelay = 500 * time.Millisecond
)

// retryTransport повторяет запрос, если стенд ответил 503
type retryTransport struct {
	next  http.RoundTripper
	delay time.Duration
}

// newRetryClient создаёт http-клиент, который переживает кратковременную недоступность стенда
func newRetryClient() *http.Client {
	return &http.Client{Transport: retryTransport{next: http.DefaultTransport, delay: retryDelay}}
}

// RoundTrip отправляет запрос и при ответе 503 повторяет его до retryAttempts раз, тело запроса при этом
// отправляется заново целиком. Остальные ответы, в том числе ошибки 4xx, возвращаются как есть
func (rt retryTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	var body []byte
	if request.Body != nil {
		var err error
		body, err = io.ReadAll(request.Body)
		_ = request.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать тело запроса для повторной отправки: %w", err)
		}
	}

	for attempt := 1; ; attempt++ {
		attemptRequest := request.Clone(request.Context())
		if body != nil {
			attemptRequest.Body = io.NopCloser(bytes.NewReader(body))
		}

		response, err := rt.next.RoundTrip(attemptRequest)
		if err != nil || response.StatusCode != http.StatusServiceUnavailable || attempt == retryAttempts {
			return response, err
		}

		_ = response.Body.Close()
		utils.LogWithLabelAndTimeStamp("Retry", fmt.Sprintf("%s %s вернул 503, попытка %d из %d",
			request.Method, request.URL.Path, attempt, retryAttempts))
		time.Sleep(rt.delay * time.Duration(attempt))
	}
}
