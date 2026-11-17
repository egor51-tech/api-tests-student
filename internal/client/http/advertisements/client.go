package advertisements

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	advModels "api-tests-template/internal/managers/advertisements/models"
)

const photosFieldName = "photos"

// HttpPostAdvertisement выполняет POST-запрос для создания объявления с фотографиями
func HttpPostAdvertisement(t *testing.T, token string, request advModels.CreateAdvertisementRequest) *http.Response {
	body, contentType := buildMultipartBody(t, request)

	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		ContentType(contentType).
		Body(body).
		Expect(t).
		End().Response
}

// HttpGetAdvertisement выполняет GET-запрос для получения объявления по id
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisements выполняет GET-запрос для поиска объявлений по названию
func HttpGetAdvertisements(t *testing.T, search string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementsPath).
		Query("search", search).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement выполняет DELETE-запрос для удаления объявления
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos выполняет GET-запрос для получения фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(fmt.Sprintf(path.AdvertisementPhotosPathTemplate, id)).
		Expect(t).
		End().Response
}

// buildMultipartBody собирает тело multipart/form-data для создания объявления и возвращает его вместе с Content-Type
func buildMultipartBody(t *testing.T, request advModels.CreateAdvertisementRequest) (string, string) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	fields := []struct{ name, value string }{
		{"title", request.Title},
		{"description", request.Description},
		{"price", request.Price},
		{"quantity", request.Quantity},
	}
	for _, field := range fields {
		if field.value == "" {
			continue
		}
		require.NoErrorf(t, writer.WriteField(field.name, field.value), "Не удалось записать поле %s в форму", field.name)
	}

	for _, photo := range request.Photos {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, photosFieldName, photo.FileName))
		header.Set("Content-Type", photo.ContentType)

		part, err := writer.CreatePart(header)
		require.NoErrorf(t, err, "Не удалось создать часть формы для файла %s", photo.FileName)

		_, err = part.Write(photo.Content)
		require.NoErrorf(t, err, "Не удалось записать файл %s в форму", photo.FileName)
	}

	require.NoError(t, writer.Close(), "Не удалось завершить multipart-форму")

	return buffer.String(), writer.FormDataContentType()
}
