package testdata

// Значения для создаваемых в тестах объявлений
const (
	AdvertisementTitlePrefix = "QA advertisement "
	AdvertisementDescription = "Описание объявления, созданного автотестом"
	AdvertisementPrice       = 1999.99
	AdvertisementQuantity    = 5

	// PriceDelta допустимая погрешность при сравнении цены
	PriceDelta = 0.001

	// MaxTitleLength максимальная допустимая длина названия объявления
	MaxTitleLength = 255
	// MaxPhotos максимальное количество фотографий в объявлении
	MaxPhotos = 3
	// PhotosOverLimit количество фотографий, превышающее допустимый максимум
	PhotosOverLimit = MaxPhotos + 1

	PhotoContentType = "image/png"
	PhotoFileNameFmt = "photo_%d.png"
)

// Значения для негативных сценариев
const (
	NegativePrice = "-1"
	// InvalidToken заведомо неверный авторизационный токен
	InvalidToken = "incorrect_token"
	// NonExistentAdvertisementID корректный по формату, но не существующий id объявления
	NonExistentAdvertisementID = "00000000-0000-4000-8000-000000000000"
)

// Параметры избранного
const (
	// FavoritesLimit максимальное количество объявлений в избранном одного пользователя
	FavoritesLimit = 50
	// FavoritesPageSize размер страницы избранного (допустимые значения по API: 10, 15, 20, 30, 50)
	FavoritesPageSize = 50
)

// Параметры генерации данных
const (
	// RandomSuffixLength длина случайной части уникальных названий и email
	RandomSuffixLength = 12
	// UserPasswordLength длина пароля случайного пользователя (минимум по API - 8)
	UserPasswordLength = 16
	UserFirstName      = "QA"
	UserEmailDomain    = "example.com"
)
