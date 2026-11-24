package favorites

import (
	"fmt"
	"net/http"

	"api-tests-template/internal/constants/messages"
	"api-tests-template/internal/constants/testdata"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// Negative: Добавляем в избранное несуществующее объявление, ожидаем 404
func (s *TestSuite) TestAddNonExistentAdvertisement() {
	base.Step(&s.Suite, "Добавляем в избранное объявление с несуществующим id", func() {
		errorResponse := favorites.AddFavoriteExpectError(s.T(), s.user.Token, testdata.NonExistentAdvertisementID, http.StatusNotFound)
		base.ResultError(http.StatusNotFound, errorResponse)

		response.RequireError(s.T(), errorResponse, messages.ErrorNotFound, messages.MessageAdvertisementNotFound)
	})
}

// Negative: Удаляем из избранного несуществующее объявление, ожидаем 404
func (s *TestSuite) TestRemoveNonExistentAdvertisement() {
	base.Step(&s.Suite, "Удаляем из избранного объявление с несуществующим id", func() {
		errorResponse := favorites.DeleteFavoriteExpectError(s.T(), s.user.Token, testdata.NonExistentAdvertisementID, http.StatusNotFound)
		base.ResultError(http.StatusNotFound, errorResponse)

		response.RequireError(s.T(), errorResponse, messages.ErrorNotFound, messages.MessageFavoriteNotFound)
	})
}

// Negative: Добавляем в избранное своё собственное объявление, ожидаем 403
func (s *TestSuite) TestAddOwnAdvertisement() {
	base.Precondition("Пользователь создаёт собственное объявление")
	ownAdvertisement, _ := s.createAdvertisement(s.user.Token)

	base.Step(&s.Suite, "Добавляем собственное объявление в избранное", func() {
		errorResponse := favorites.AddFavoriteExpectError(s.T(), s.user.Token, ownAdvertisement.Id, http.StatusForbidden)
		base.ResultError(http.StatusForbidden, errorResponse)

		response.RequireError(s.T(), errorResponse, messages.ErrorForbidden, messages.MessageOwnAdvertisement)
	})
}

// Negative: Повторно добавляем в избранное то же объявление, ожидаем 409
func (s *TestSuite) TestAddAdvertisementTwice() {
	base.Precondition("Второй пользователь создаёт объявление")
	advertisement, _ := s.createAdvertisement(s.seller.Token)

	base.Step(&s.Suite, "Добавляем объявление в избранное первый раз", func() {
		favorites.AddFavorite(s.T(), s.user.Token, advertisement.Id)
	})
	base.RemoveFavoriteAfterTest(s.T(), s.user.Token, advertisement.Id)

	base.Step(&s.Suite, "Добавляем то же объявление в избранное повторно", func() {
		errorResponse := favorites.AddFavoriteExpectError(s.T(), s.user.Token, advertisement.Id, http.StatusConflict)
		base.ResultError(http.StatusConflict, errorResponse)

		response.RequireError(s.T(), errorResponse, messages.ErrorAlreadyExists, messages.MessageAlreadyInFavorites)
	})
}

// Negative: Добавляем в избранное без id объявления, ожидаем 400
func (s *TestSuite) TestAddWithoutAdvertisementId() {
	base.Step(&s.Suite, "Добавляем в избранное запись с пустым advertisement_id", func() {
		errorResponse := favorites.AddFavoriteExpectError(s.T(), s.user.Token, "", http.StatusBadRequest)
		base.ResultError(http.StatusBadRequest, errorResponse)

		response.RequireError(s.T(), errorResponse, messages.ErrorInvalidRequest, messages.MessageInvalidUUIDEmpty)
	})
}

// Negative: Обращаемся к ручкам избранного без токена и с неверным токеном, ожидаем 401
func (s *TestSuite) TestFavoritesUnauthorized() {
	endpoints := []struct {
		name string
		call func(token string) response.ErrorResponse
	}{
		{
			name: "GET /favorites",
			call: func(token string) response.ErrorResponse {
				return favorites.GetFavoritesExpectError(s.T(), token, http.StatusUnauthorized)
			},
		},
		{
			name: "POST /favorites",
			call: func(token string) response.ErrorResponse {
				return favorites.AddFavoriteExpectError(s.T(), token, testdata.NonExistentAdvertisementID, http.StatusUnauthorized)
			},
		},
		{
			name: "DELETE /favorites/{ad_id}",
			call: func(token string) response.ErrorResponse {
				return favorites.DeleteFavoriteExpectError(s.T(), token, testdata.NonExistentAdvertisementID, http.StatusUnauthorized)
			},
		},
	}
	tokens := []struct {
		name  string
		token string
	}{
		{name: "без токена", token: ""},
		{name: "неверный токен", token: testdata.InvalidToken},
	}

	for _, endpoint := range endpoints {
		for _, token := range tokens {
			base.Case(&s.Suite, fmt.Sprintf("%s, %s", endpoint.name, token.name), func() {
				errorResponse := endpoint.call(token.token)
				base.ResultError(http.StatusUnauthorized, errorResponse)

				base.RequireUnauthorizedError(s.T(), errorResponse, token.token)
			})
		}
	}
}
