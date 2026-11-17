package models

// CreateAdvertisementRequest поля формы создания объявления (multipart/form-data).
// Пустое строковое поле в запрос не попадает, так можно проверять отсутствие параметра.
type CreateAdvertisementRequest struct {
	Title       string
	Description string
	Price       string
	Quantity    string
	Photos      []PhotoFile
}

// PhotoFile файл фотографии для загрузки
type PhotoFile struct {
	FileName    string
	ContentType string
	Content     []byte
}
