package createAdvertisement

import (
	"net/http"
	"strings"

	"api-tests-template/internal/constants/testdata"
	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"

	base "api-tests-template/tests"
)

// Negative: Пытаемся создать объявление без обязательных или с некорректными параметрами, ожидаем 400
func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	testCases := []struct {
		name   string
		modify func(request *advModels.CreateAdvertisementRequest)
	}{
		{
			name:   "Пустое название",
			modify: func(request *advModels.CreateAdvertisementRequest) { request.Title = "" },
		},
		{
			name: "Слишком длинное название",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Title = strings.Repeat("a", testdata.MaxTitleLength+1)
			},
		},
		{
			name:   "Пустое описание",
			modify: func(request *advModels.CreateAdvertisementRequest) { request.Description = "" },
		},
		{
			name:   "Отрицательная цена",
			modify: func(request *advModels.CreateAdvertisementRequest) { request.Price = testdata.NegativePrice },
		},
		{
			name:   "Нет фотографий",
			modify: func(request *advModels.CreateAdvertisementRequest) { request.Photos = nil },
		},
		{
			name: "Фотографий больше допустимого",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				*request = advertisements.NewAdvertisementRequest(testdata.PhotosOverLimit)
			},
		},
	}

	for _, testCase := range testCases {
		base.Case(&s.Suite, testCase.name, func() {
			request := advertisements.NewAdvertisementRequest(1)
			testCase.modify(&request)

			errorResponse := advertisements.CreateAdvertisementExpectError(s.T(), s.loginData.Token, request, http.StatusBadRequest)
			base.ResultError(http.StatusBadRequest, errorResponse)
		})
	}
}

// Negative: Пытаемся создать объявление без токена и с неверным токеном, ожидаем 401
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	testCases := []struct {
		name  string
		token string
	}{
		{name: "Без токена", token: ""},
		{name: "Неверный токен", token: testdata.InvalidToken},
	}

	for _, testCase := range testCases {
		base.Case(&s.Suite, testCase.name, func() {
			request := advertisements.NewAdvertisementRequest(1)

			errorResponse := advertisements.CreateAdvertisementExpectError(s.T(), testCase.token, request, http.StatusUnauthorized)
			base.RequireUnauthorizedError(s.T(), errorResponse, testCase.token)
			base.ResultError(http.StatusUnauthorized, errorResponse)
		})
	}
}
