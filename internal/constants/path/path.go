package path

const (
	AuthLoginPath        = "/auth/login"
	AuthRegisterPath     = "/auth/register"
	MyAdvertisementsPath = "/my/advertisements"

	AdvertisementPath  = "/advertisement"
	AdvertisementsPath = "/advertisements"
	// AdvertisementPhotosPathTemplate шаблон пути фотографий объявления, подставляется id объявления
	AdvertisementPhotosPathTemplate = "/advertisements/%s/photos"

	FavoritesPath     = "/favorites"
	FavoritesBulkPath = "/favorites/bulk"
	// FavoritePathTemplate шаблон пути конкретного избранного, подставляется id объявления
	FavoritePathTemplate = "/favorites/%s"
)
