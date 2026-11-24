package favorites

import (
	"net/http"
	"slices"
	"testing"

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

// ContainsAdvertisement проверяет, есть ли объявление с указанным id в списке избранного
func ContainsAdvertisement(list models.FavoriteListResponse, adId string) bool {
	return slices.ContainsFunc(list.Items, func(item models.FavoriteResponse) bool {
		return item.Advertisement.Id == adId
	})
}
