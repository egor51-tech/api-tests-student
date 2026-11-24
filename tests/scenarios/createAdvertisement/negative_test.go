package createAdvertisement

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/messages"
	"api-tests-template/internal/constants/testdata"
	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"

	base "api-tests-template/tests"
)

// Negative: Пытаемся создать объявление без обязательных или с некорректными параметрами, ожидаем 400
// и свой текст ошибки для каждого случая
func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	testCases := []struct {
		name   string
		modify func(request *advModels.CreateAdvertisementRequest)
		// message ожидаемый текст ошибки; если пустой, проверяются только status code, error и непустой message
		message string
	}{
		{
			name:    "Пустое название",
			modify:  func(request *advModels.CreateAdvertisementRequest) { request.Title = "" },
			message: messages.MessageTitleRequired,
		},
		{
			name: "Слишком длинное название",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Title = strings.Repeat("a", testdata.MaxTitleLength+1)
			},
			message: fmt.Sprintf(messages.MessageTitleTooLongFmt, testdata.MaxTitleLength, testdata.MaxTitleLength+1),
		},
		{
			name:    "Пустое описание",
			modify:  func(request *advModels.CreateAdvertisementRequest) { request.Description = "" },
			message: messages.MessageDescriptionRequired,
		},
		{
			name:    "Отрицательная цена",
			modify:  func(request *advModels.CreateAdvertisementRequest) { request.Price = testdata.NegativePrice },
			message: messages.MessagePriceOutOfRange,
		},
		{
			name:   "Цена не число",
			modify: func(request *advModels.CreateAdvertisementRequest) { request.Price = testdata.NotNumber },
		},
		{
			name:   "Количество не число",
			modify: func(request *advModels.CreateAdvertisementRequest) { request.Quantity = testdata.NotNumber },
		},
		{
			name:    "Нет фотографий",
			modify:  func(request *advModels.CreateAdvertisementRequest) { request.Photos = nil },
			message: messages.MessagePhotoRequired,
		},
		{
			name: "Фотографий больше допустимого",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Photos = advertisements.LoadPhotos(s.T(), testdata.PhotosOverLimit)
			},
			message: messages.MessagePhotosOverLimit,
		},
		{
			name: "Неподдерживаемый формат файла",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Photos = []advModels.PhotoFile{
					advertisements.LoadPhoto(s.T(), testdata.NotImageFileName, testdata.NotImageContentType),
				}
			},
		},
	}

	for _, testCase := range testCases {
		base.Case(&s.Suite, testCase.name, func() {
			request := advertisements.NewAdvertisementRequest(advertisements.LoadPhotos(s.T(), 1))
			testCase.modify(&request)

			errorResponse := advertisements.CreateAdvertisementExpectError(s.T(), s.loginData.Token, request, http.StatusBadRequest)
			base.ResultError(http.StatusBadRequest, errorResponse)

			require.Equal(s.T(), messages.ErrorInvalidRequest, errorResponse.Error, "Некорректное значение error параметра")
			if testCase.message != "" {
				require.Equal(s.T(), testCase.message, errorResponse.Message, "Некорректное значение message параметра")
			}
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
			request := advertisements.NewAdvertisementRequest(advertisements.LoadPhotos(s.T(), 1))

			errorResponse := advertisements.CreateAdvertisementExpectError(s.T(), testCase.token, request, http.StatusUnauthorized)
			base.ResultError(http.StatusUnauthorized, errorResponse)

			base.RequireUnauthorizedError(s.T(), errorResponse, testCase.token)
		})
	}
}
