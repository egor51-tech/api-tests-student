package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/constants/messages"
	"api-tests-template/internal/constants/testdata"
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

// DeleteAdvertisementAfterTest регистрирует удаление объявления после завершения теста.
// Вызывать нужно на уровне самого теста, а не внутри Step: иначе объявление удалится сразу после шага.
func DeleteAdvertisementAfterTest(t *testing.T, token, adId string) {
	t.Cleanup(func() {
		utils.LogWithLabelAndTimeStamp("Cleanup", "Удаляем объявление "+adId)
		advertisements.CleanupAdvertisement(t, token, adId)
	})
}

// RemoveFavoriteAfterTest регистрирует удаление объявления из избранного после завершения теста
func RemoveFavoriteAfterTest(t *testing.T, token, adId string) {
	t.Cleanup(func() {
		utils.LogWithLabelAndTimeStamp("Cleanup", "Убираем из избранного объявление "+adId)
		favorites.CleanupFavorite(t, token, adId)
	})
}

// RemoveFavoritesAfterTest регистрирует удаление нескольких объявлений из избранного после завершения теста
func RemoveFavoritesAfterTest(t *testing.T, token string, adIds []string) {
	t.Cleanup(func() {
		utils.LogWithLabelAndTimeStamp("Cleanup", fmt.Sprintf("Убираем из избранного объявления: %d шт.", len(adIds)))
		for _, adId := range adIds {
			favorites.CleanupFavorite(t, token, adId)
		}
	})
}

// RequireUnauthorizedError проверяет тело ответа 401: значение error, а для неверного токена ещё и текст message
func RequireUnauthorizedError(t *testing.T, errorResponse response.ErrorResponse, token string) {
	require.Equal(t, messages.ErrorUnauthorized, errorResponse.Error, "Некорректное значение error параметра")

	if token == testdata.InvalidToken {
		require.Equal(t, messages.MessageInvalidOrExpired, errorResponse.Message, "Некорректное значение message параметра")
	}
}
