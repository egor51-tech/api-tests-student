package favorites

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	favModels "api-tests-template/internal/managers/favorites/models"
)

// HttpGetFavorites выполняет GET-запрос для получения избранного; поиск по названию применяется, если он задан
func HttpGetFavorites(t *testing.T, token, search string, pageSize int) *http.Response {
	request := apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page_size", strconv.Itoa(pageSize))
	if search != "" {
		request = request.Query("search", search)
	}

	return request.Expect(t).End().Response
}

// HttpPostFavorite выполняет POST-запрос для добавления объявления в избранное
func HttpPostFavorite(t *testing.T, token string, request favModels.AddFavoriteRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite выполняет DELETE-запрос для удаления объявления из избранного
func HttpDeleteFavorite(t *testing.T, token, adId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(fmt.Sprintf(path.FavoritePathTemplate, adId)).
		Expect(t).
		End().Response
}
