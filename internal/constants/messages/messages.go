package messages

// Значения поля error в ответах API
const (
	ErrorUnauthorized   = "unauthorized"
	ErrorInvalidRequest = "invalid_request"
	ErrorNotFound       = "not_found"
	ErrorForbidden      = "forbidden"
	ErrorAlreadyExists  = "already_exists"
)

// Тексты поля message в ответах 401
const (
	MessageAuthHeaderRequired = "Authorization header required"
	MessageInvalidOrExpired   = "Invalid or expired token"
)

// Тексты поля message при создании объявления
const (
	MessageTitleRequired       = "title is required"
	MessageTitleTooLongFmt     = "title must not exceed %d characters provided %d"
	MessageDescriptionRequired = "description is required"
	MessagePriceOutOfRange     = "price must be between 0 and 1000000 inclusive"
	MessagePhotoRequired       = "at least 1 photo is required"
	MessagePhotosOverLimit     = "maximum 3 photos allowed"
)

// Тексты поля message в ответах избранного
const (
	MessageAdvertisementNotFound = "Advertisement not found"
	MessageFavoriteNotFound      = "Favorite entry not found"
	MessageOwnAdvertisement      = "Cannot add own advertisement to favorites"
	MessageAlreadyInFavorites    = "Advertisement is already in favorites"
	MessageInvalidUUIDEmpty      = "invalid UUID length: 0"
)
