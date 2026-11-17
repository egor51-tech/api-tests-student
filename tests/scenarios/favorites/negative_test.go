package favorites

import (
	"fmt"
	"net/http"

	"api-tests-template/internal/constants/testdata"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// Negative: Добавляем в избранное несуществующее объявление, ожидаем 404
func (s *TestSuite) TestAddNonExistentAdvertisement() {
	base.Step(&s.Suite, "Добавляем в избранное объявление с несуществующим id", func() {
		errorResponse := favorites.AddFavoriteExpectError(s.T(), s.user.Token, testdata.NonExistentAdvertisementID, http.StatusNotFound)
		base.ResultError(http.StatusNotFound, errorResponse)
	})
}

// Negative: Удаляем из избранного несуществующее объявление, ожидаем 404
func (s *TestSuite) TestRemoveNonExistentAdvertisement() {
	base.Step(&s.Suite, "Удаляем из избранного объявление с несуществующим id", func() {
		errorResponse := favorites.DeleteFavoriteExpectError(s.T(), s.user.Token, testdata.NonExistentAdvertisementID, http.StatusNotFound)
		base.ResultError(http.StatusNotFound, errorResponse)
	})
}

// Negative: Добавляем в избранное своё собственное объявление, ожидаем 403
func (s *TestSuite) TestAddOwnAdvertisement() {
	base.Precondition("Пользователь создаёт собственное объявление")
	ownAdvertisement := advertisements.CreateAdvertisement(s.T(), s.user.Token, advertisements.NewAdvertisementRequest(1))
	base.DeleteAdvertisementAfterTest(s.T(), s.user.Token, ownAdvertisement.Id)

	base.Step(&s.Suite, "Добавляем собственное объявление в избранное", func() {
		errorResponse := favorites.AddFavoriteExpectError(s.T(), s.user.Token, ownAdvertisement.Id, http.StatusForbidden)
		base.ResultError(http.StatusForbidden, errorResponse)
	})
}

// Negative: Повторно добавляем в избранное то же объявление, ожидаем 409
func (s *TestSuite) TestAddAdvertisementTwice() {
	base.Precondition("Второй пользователь создаёт объявление")
	advertisement := s.createSellerAdvertisement()
	base.RemoveFavoriteAfterTest(s.T(), s.user.Token, advertisement.Id)

	base.Step(&s.Suite, "Добавляем объявление в избранное первый раз", func() {
		favorites.AddFavorite(s.T(), s.user.Token, advertisement.Id)
	})

	base.Step(&s.Suite, "Добавляем то же объявление в избранное повторно", func() {
		errorResponse := favorites.AddFavoriteExpectError(s.T(), s.user.Token, advertisement.Id, http.StatusConflict)
		base.ResultError(http.StatusConflict, errorResponse)
	})
}

// Negative: Добавляем в избранное без id объявления, ожидаем 400
func (s *TestSuite) TestAddWithoutAdvertisementId() {
	base.Step(&s.Suite, "Добавляем в избранное запись с пустым advertisement_id", func() {
		errorResponse := favorites.AddFavoriteExpectError(s.T(), s.user.Token, "", http.StatusBadRequest)
		base.ResultError(http.StatusBadRequest, errorResponse)
	})
}

// Negative: Добавляем объявление в избранное, в котором уже лежит максимум записей, ожидаем 400
func (s *TestSuite) TestAddOverFavoritesLimit() {
	base.Precondition(fmt.Sprintf("Второй пользователь создаёт %d объявлений для заполнения избранного и ещё одно сверх лимита",
		testdata.FavoritesLimit))
	adIds := make([]string, 0, testdata.FavoritesLimit)
	for range testdata.FavoritesLimit {
		adIds = append(adIds, s.createSellerAdvertisement().Id)
	}
	extraAdvertisement := s.createSellerAdvertisement()
	base.RemoveFavoritesAfterTest(s.T(), s.user.Token, adIds)

	base.Step(&s.Suite, "Заполняем избранное до лимита одним запросом", func() {
		result := favorites.AddFavoritesBulk(s.T(), s.user.Token, adIds)
		base.Result(fmt.Sprintf("Добавлено %d, пропущено %d, лимит достигнут: %t", len(result.Added), len(result.Skipped), result.LimitReached))
	})

	base.Step(&s.Suite, "Добавляем ещё одно объявление сверх лимита", func() {
		errorResponse := favorites.AddFavoriteExpectError(s.T(), s.user.Token, extraAdvertisement.Id, http.StatusBadRequest)
		base.ResultError(http.StatusBadRequest, errorResponse)
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

				base.RequireUnauthorizedError(s.T(), errorResponse, token.token)
				base.ResultError(http.StatusUnauthorized, errorResponse)
			})
		}
	}
}
