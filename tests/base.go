package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/constants/messages"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
)

// SetupSuite действия перед запуском сьюта с тестами
func SetupSuite() {
	log.Println("Init environment variables")
	utils.LoadEnv()
}

// TearDownSuite действия после запуска сьюта с тестами
func TearDownSuite() {
	log.Println("Tear down suite")
}

// Precondition выводит текст в консоль с меткой [Precondition]
func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Precondition", text)
}

// Result выводит результат проверки в консоль с меткой [Result]
func Result(text string) {
	utils.LogWithLabelAndTimeStamp("Result", text)
}

// ResultError выводит в консоль полученную от API ошибку с меткой [Result]
func ResultError(statusCode int, errorResponse response.ErrorResponse) {
	Result(fmt.Sprintf("Получен %d, error=%q, message=%q", statusCode, errorResponse.Error, errorResponse.Message))
}

// Step выполняет шаг теста как подтест и пишет его название в консоль с меткой [Step].
// Если шаг упал, тест останавливается: следующие шаги зависят от результата предыдущих.
func Step(s *suite.Suite, name string, step func()) {
	utils.LogWithLabelAndTimeStamp("Step", name)
	if !s.Run(name, step) {
		s.T().FailNow()
	}
}

// Case выполняет один вариант параметризованного теста как подтест. В отличие от Step,
// падение одного варианта не мешает проверке остальных.
func Case(s *suite.Suite, name string, testCase func()) {
	utils.LogWithLabelAndTimeStamp("Case", name)
	s.Run(name, testCase)
}

// DeleteAdvertisementAfterTest регистрирует удаление объявления после завершения теста и возвращает
// функцию отмены. Если тест сам удалил объявление, нужно вызвать отмену, чтобы очистка не удаляла его повторно.
// Регистрировать нужно на уровне самого теста, а не внутри Step: иначе объявление удалится сразу после шага.
func DeleteAdvertisementAfterTest(t *testing.T, token, adId string) (cancel func()) {
	cancelled := false
	t.Cleanup(func() {
		if cancelled {
			return
		}
		utils.LogWithLabelAndTimeStamp("Cleanup", "Удаляем объявление "+adId)
		advertisements.DeleteAdvertisement(t, token, adId)
	})

	return func() { cancelled = true }
}

// RemoveFavoriteAfterTest регистрирует удаление объявления из избранного после завершения теста и возвращает
// функцию отмены. Если тест сам убрал объявление из избранного, нужно вызвать отмену.
func RemoveFavoriteAfterTest(t *testing.T, token, adId string) (cancel func()) {
	cancelled := false
	t.Cleanup(func() {
		if cancelled {
			return
		}
		utils.LogWithLabelAndTimeStamp("Cleanup", "Убираем из избранного объявление "+adId)
		favorites.DeleteFavorite(t, token, adId)
	})

	return func() { cancelled = true }
}

// RequireUnauthorizedError проверяет тело ответа 401: для запроса без токена и с неверным токеном тексты разные
func RequireUnauthorizedError(t *testing.T, errorResponse response.ErrorResponse, token string) {
	t.Helper()

	expectedMessage := messages.MessageInvalidOrExpired
	if token == "" {
		expectedMessage = messages.MessageAuthHeaderRequired
	}

	response.RequireError(t, errorResponse, messages.ErrorUnauthorized, expectedMessage)
}
