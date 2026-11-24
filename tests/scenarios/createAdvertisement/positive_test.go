package createAdvertisement

import (
	"github.com/stretchr/testify/require"

	"api-tests-template/internal/constants/testdata"
	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"

	base "api-tests-template/tests"
)

// Positive: Создаём объявление со всеми параметрами и проверяем его через карточку, поиск и ручку фотографий
func (s *TestSuite) TestCreateAdvertisementPositive() {
	token := s.loginData.Token

	var request advModels.CreateAdvertisementRequest
	var created advModels.AdvertisementResponse

	base.Step(&s.Suite, "Создаём объявление с названием, описанием, ценой, количеством и тремя фотографиями", func() {
		request = advertisements.NewAdvertisementRequest(advertisements.LoadPhotos(s.T(), testdata.MaxPhotos))
		created = advertisements.CreateAdvertisement(s.T(), token, request)

		require.Len(s.T(), created.Photos, testdata.MaxPhotos, "Количество фотографий в ответе на создание не совпадает")
		base.Result("Создано объявление с id " + created.Id)
	})
	base.DeleteAdvertisementAfterTest(s.T(), token, created.Id)

	base.Step(&s.Suite, "Получаем объявление по id и сравниваем параметры с переданными при создании", func() {
		advertisement := advertisements.GetAdvertisement(s.T(), created.Id)

		require.Equal(s.T(), request.Title, advertisement.Title, "Название не совпадает")
		require.Equal(s.T(), request.Description, advertisement.Description, "Описание не совпадает")
		require.InDelta(s.T(), testdata.AdvertisementPrice, advertisement.Price, testdata.PriceDelta, "Цена не совпадает")
		require.Equal(s.T(), testdata.AdvertisementQuantity, advertisement.Quantity, "Количество не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, advertisement.UserId, "Автор объявления не совпадает с авторизованным пользователем")
		require.ElementsMatch(s.T(), advertisements.PhotoIds(created.Photos), advertisements.PhotoIds(advertisement.Photos),
			"Фотографии объявления не совпадают с загруженными")
		base.Result("Параметры объявления совпадают с переданными")
	})

	base.Step(&s.Suite, "Ищем объявление по полному названию", func() {
		found := advertisements.SearchAdvertisements(s.T(), request.Title)

		require.Truef(s.T(), advertisements.ContainsAdvertisement(found, created.Id),
			"Объявление %s не найдено в поиске по названию %q", created.Id, request.Title)
		base.Result("Объявление найдено в поиске")
	})

	base.Step(&s.Suite, "Получаем фотографии объявления", func() {
		photos := advertisements.GetAdvertisementPhotos(s.T(), token, created.Id)

		require.Len(s.T(), photos, testdata.MaxPhotos, "Количество загруженных фотографий не совпадает")
		require.ElementsMatch(s.T(), advertisements.PhotoIds(created.Photos), advertisements.PhotoIds(photos),
			"Id фотографий не совпадают с полученными при создании")
		for _, photo := range photos {
			require.NotEmptyf(s.T(), photo.Url, "У фотографии %s пустая ссылка", photo.Id)
		}
		base.Result("Все фотографии доступны")
	})
}
