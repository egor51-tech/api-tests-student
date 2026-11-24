package favorites

import (
	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

// Positive: Полный цикл работы с избранным: добавляем объявление, видим его в списке, удаляем и видим, что оно пропало
func (s *TestSuite) TestFavoritesFullCycle() {
	base.Precondition("Второй пользователь создаёт объявление")
	advertisement, _ := s.createAdvertisement(s.seller.Token)

	base.Step(&s.Suite, "Добавляем объявление в избранное", func() {
		favorite := favorites.AddFavorite(s.T(), s.user.Token, advertisement.Id)

		require.NotEmpty(s.T(), favorite.FavoriteId, "В ответе на добавление отсутствует favorite_id")
		require.Equal(s.T(), advertisement.Id, favorite.Advertisement.Id, "В ответе на добавление вернулось другое объявление")
		base.Result("Объявление добавлено в избранное, favorite_id " + favorite.FavoriteId)
	})
	cancelFavoriteCleanup := base.RemoveFavoriteAfterTest(s.T(), s.user.Token, advertisement.Id)

	base.Step(&s.Suite, "Проверяем, что объявление появилось в списке избранного", func() {
		list := favorites.GetFavorites(s.T(), s.user.Token, advertisement.Title)

		require.Truef(s.T(), favorites.ContainsAdvertisement(list, advertisement.Id),
			"Объявление %s не найдено в избранном", advertisement.Id)
		base.Result("Объявление есть в списке избранного")
	})

	base.Step(&s.Suite, "Удаляем объявление из избранного", func() {
		favorites.DeleteFavorite(s.T(), s.user.Token, advertisement.Id)
		base.Result("Объявление удалено из избранного")
	})
	cancelFavoriteCleanup()

	base.Step(&s.Suite, "Проверяем, что объявление пропало из списка избранного", func() {
		list := favorites.GetFavorites(s.T(), s.user.Token, advertisement.Title)

		require.Falsef(s.T(), favorites.ContainsAdvertisement(list, advertisement.Id),
			"Объявление %s осталось в избранном после удаления", advertisement.Id)
		base.Result("Объявления нет в списке избранного")
	})
}

// Positive: После удаления самого объявления оно исчезает из избранного пользователя
func (s *TestSuite) TestDeleteAdvertisementRemovesItFromFavorites() {
	base.Precondition("Второй пользователь создаёт объявление")
	advertisement, cancelAdvertisementCleanup := s.createAdvertisement(s.seller.Token)

	base.Step(&s.Suite, "Добавляем объявление в избранное", func() {
		favorites.AddFavorite(s.T(), s.user.Token, advertisement.Id)

		list := favorites.GetFavorites(s.T(), s.user.Token, advertisement.Title)
		require.Truef(s.T(), favorites.ContainsAdvertisement(list, advertisement.Id),
			"Объявление %s не найдено в избранном", advertisement.Id)
		base.Result("Объявление есть в списке избранного")
	})

	base.Step(&s.Suite, "Владелец удаляет объявление", func() {
		advertisements.DeleteAdvertisement(s.T(), s.seller.Token, advertisement.Id)
		base.Result("Объявление удалено")
	})
	// Объявление уже удалено, повторно в очистке его удалять не нужно
	cancelAdvertisementCleanup()

	base.Step(&s.Suite, "Проверяем, что объявление пропало из избранного", func() {
		list := favorites.GetFavorites(s.T(), s.user.Token, advertisement.Title)

		require.Falsef(s.T(), favorites.ContainsAdvertisement(list, advertisement.Id),
			"Удалённое объявление %s осталось в избранном", advertisement.Id)
		base.Result("Удалённого объявления нет в списке избранного")
	})
}
