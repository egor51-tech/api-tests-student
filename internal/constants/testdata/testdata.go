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
	MaxTitleLength = 50
	// MaxPhotos максимальное количество фотографий в объявлении
	MaxPhotos = 3
	// PhotosOverLimit количество фотографий, превышающее допустимый максимум
	PhotosOverLimit = MaxPhotos + 1
)

// Файлы для загрузки лежат в папке testdata рядом с тестами
const (
	// PhotosDir папка с файлами относительно пакета теста
	PhotosDir = "testdata"
	// PhotoFileNameFmt шаблон имени корректной фотографии, подставляется номер от 1
	PhotoFileNameFmt = "photo_%d.png"
	PhotoContentType = "image/png"

	// NotImageFileName текстовый файл, который нельзя загружать как фотографию
	NotImageFileName    = "not_an_image.txt"
	NotImageContentType = "text/plain"
)

// Значения для негативных сценариев
const (
	NegativePrice = "-1"
	// NotNumber строка, которая не является числом
	NotNumber = "abc"
	// InvalidToken заведомо неверный авторизационный токен
	InvalidToken = "incorrect_token"
	// NonExistentAdvertisementID корректный по формату, но не существующий id объявления
	NonExistentAdvertisementID = "00000000-0000-4000-8000-000000000000"
)

// Параметры избранного
const (
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
