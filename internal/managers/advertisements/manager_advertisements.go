package advertisements

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/constants/testdata"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

// NewAdvertisementRequest возвращает данные корректного объявления с уникальным названием и нужным числом фотографий
func NewAdvertisementRequest(photosCount int) models.CreateAdvertisementRequest {
	photos := make([]models.PhotoFile, 0, photosCount)
	for i := range photosCount {
		photos = append(photos, models.PhotoFile{
			FileName:    fmt.Sprintf(testdata.PhotoFileNameFmt, i+1),
			ContentType: testdata.PhotoContentType,
			Content:     utils.GeneratePNG(i),
		})
	}

	return models.CreateAdvertisementRequest{
		Title:       testdata.AdvertisementTitlePrefix + utils.RandomString(testdata.RandomSuffixLength),
		Description: testdata.AdvertisementDescription,
		Price:       strconv.FormatFloat(testdata.AdvertisementPrice, 'f', 2, 64),
		Quantity:    strconv.Itoa(testdata.AdvertisementQuantity),
		Photos:      photos,
	}
}

// CreateAdvertisement создаёт объявление и возвращает его из ответа, ожидает 201
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.AdvertisementResponse {
	rawResponse := advertisements.HttpPostAdvertisement(t, token, request)
	body := response.ReadBodyWithStatus(t, rawResponse, http.StatusCreated)

	created := response.Decode[models.AdvertisementResponse](t, body)
	require.NotEmpty(t, created.Id, "В ответе на создание объявления отсутствует id")

	return created
}

// CreateAdvertisementExpectError пытается создать объявление и проверяет статус ошибки и тело ответа
func CreateAdvertisementExpectError(t *testing.T, token string, request models.CreateAdvertisementRequest,
	expectedStatusCode int) response.ErrorResponse {
	rawResponse := advertisements.HttpPostAdvertisement(t, token, request)
	body := response.ReadBodyWithStatus(t, rawResponse, expectedStatusCode)

	return response.DecodeError(t, body)
}

// GetAdvertisement возвращает объявление по id, ожидает 200
func GetAdvertisement(t *testing.T, id string) models.AdvertisementResponse {
	rawResponse := advertisements.HttpGetAdvertisement(t, id)
	body := response.ReadBodyWithStatus(t, rawResponse, http.StatusOK)

	return response.Decode[models.AdvertisementResponse](t, body)
}

// SearchAdvertisements ищет объявления по названию, ожидает 200
func SearchAdvertisements(t *testing.T, title string) models.AdvertisementListResponse {
	rawResponse := advertisements.HttpGetAdvertisements(t, title)
	body := response.ReadBodyWithStatus(t, rawResponse, http.StatusOK)

	return response.Decode[models.AdvertisementListResponse](t, body)
}

// GetAdvertisementPhotos возвращает фотографии объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.PhotoResponse {
	rawResponse := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	body := response.ReadBodyWithStatus(t, rawResponse, http.StatusOK)

	return response.Decode[[]models.PhotoResponse](t, body)
}

// DeleteAdvertisement удаляет объявление, ожидает 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	rawResponse := advertisements.HttpDeleteAdvertisement(t, token, id)
	response.ReadBodyWithStatus(t, rawResponse, http.StatusNoContent)
}

// CleanupAdvertisement удаляет объявление после теста. Ответ 404 допустим:
// сам тест мог уже удалить это объявление.
func CleanupAdvertisement(t *testing.T, token, id string) {
	rawResponse := advertisements.HttpDeleteAdvertisement(t, token, id)
	body := response.ReadBody(t, rawResponse)

	require.Containsf(t, []int{http.StatusNoContent, http.StatusNotFound}, rawResponse.StatusCode,
		"Не удалось удалить объявление %s после теста. Тело ответа: %s", id, string(body))
}

// ContainsAdvertisement проверяет, есть ли объявление с указанным id в результатах поиска
func ContainsAdvertisement(list models.AdvertisementListResponse, id string) bool {
	return slices.ContainsFunc(list.Items, func(item models.AdvertisementResponse) bool {
		return item.Id == id
	})
}

// PhotoIds возвращает id фотографий в том порядке, в котором они пришли
func PhotoIds(photos []models.PhotoResponse) []string {
	ids := make([]string, 0, len(photos))
	for _, photo := range photos {
		ids = append(ids, photo.Id)
	}

	return ids
}
