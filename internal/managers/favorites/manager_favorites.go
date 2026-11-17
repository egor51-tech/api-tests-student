package favorites

import (
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/constants/testdata"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет объявление в избранное, ожидает 201
func AddFavorite(t *testing.T, token, adId string) models.FavoriteResponse {
	rawResponse := favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementId: adId})
	body := response.ReadBodyWithStatus(t, rawResponse, http.StatusCreated)

	return response.Decode[models.FavoriteResponse](t, body)
}

// AddFavoriteExpectError пытается добавить объявление в избранное и проверяет статус ошибки и тело ответа
func AddFavoriteExpectError(t *testing.T, token, adId string, expectedStatusCode int) response.ErrorResponse {
	rawResponse := favorites.HttpPostFavorite(t, token, models.AddFavoriteRequest{AdvertisementId: adId})
	body := response.ReadBodyWithStatus(t, rawResponse, expectedStatusCode)

	return response.DecodeError(t, body)
}

// AddFavoritesBulk добавляет несколько объявлений в избранное одним запросом, ожидает 200
func AddFavoritesBulk(t *testing.T, token string, adIds []string) models.BulkAddFavoritesResponse {
	rawResponse := favorites.HttpPostFavoritesBulk(t, token, models.BulkAddFavoritesRequest{AdvertisementIds: adIds})
	body := response.ReadBodyWithStatus(t, rawResponse, http.StatusOK)

	return response.Decode[models.BulkAddFavoritesResponse](t, body)
}

// GetFavorites возвращает избранное пользователя, отфильтрованное по названию (пустая строка - без фильтра), ожидает 200
func GetFavorites(t *testing.T, token, search string) models.FavoriteListResponse {
	rawResponse := favorites.HttpGetFavorites(t, token, search, testdata.FavoritesPageSize)
	body := response.ReadBodyWithStatus(t, rawResponse, http.StatusOK)

	return response.Decode[models.FavoriteListResponse](t, body)
}

// GetFavoritesExpectError запрашивает избранное и проверяет статус ошибки и тело ответа
func GetFavoritesExpectError(t *testing.T, token string, expectedStatusCode int) response.ErrorResponse {
	rawResponse := favorites.HttpGetFavorites(t, token, "", testdata.FavoritesPageSize)
	body := response.ReadBodyWithStatus(t, rawResponse, expectedStatusCode)

	return response.DecodeError(t, body)
}

// DeleteFavorite удаляет объявление из избранного, ожидает 204
func DeleteFavorite(t *testing.T, token, adId string) {
	rawResponse := favorites.HttpDeleteFavorite(t, token, adId)
	response.ReadBodyWithStatus(t, rawResponse, http.StatusNoContent)
}

// DeleteFavoriteExpectError пытается удалить объявление из избранного и проверяет статус ошибки и тело ответа
func DeleteFavoriteExpectError(t *testing.T, token, adId string, expectedStatusCode int) response.ErrorResponse {
	rawResponse := favorites.HttpDeleteFavorite(t, token, adId)
	body := response.ReadBodyWithStatus(t, rawResponse, expectedStatusCode)

	return response.DecodeError(t, body)
}

// CleanupFavorite убирает объявление из избранного после теста. Ответ 404 допустим:
// сам тест мог уже удалить запись или само объявление.
func CleanupFavorite(t *testing.T, token, adId string) {
	rawResponse := favorites.HttpDeleteFavorite(t, token, adId)
	body := response.ReadBody(t, rawResponse)

	require.Containsf(t, []int{http.StatusNoContent, http.StatusNotFound}, rawResponse.StatusCode,
		"Не удалось убрать объявление %s из избранного после теста. Тело ответа: %s", adId, string(body))
}

// ContainsAdvertisement проверяет, есть ли объявление с указанным id в списке избранного
func ContainsAdvertisement(list models.FavoriteListResponse, adId string) bool {
	return slices.ContainsFunc(list.Items, func(item models.FavoriteResponse) bool {
		return item.Advertisement.Id == adId
	})
}
