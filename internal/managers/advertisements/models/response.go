package models

type PhotoResponse struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}

// AdvertisementResponse объявление, в таком виде оно приходит при создании, получении по id и в списках
type AdvertisementResponse struct {
	Id          string          `json:"id"`
	UserId      string          `json:"user_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	Photos      []PhotoResponse `json:"photos"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

type AdvertisementListResponse struct {
	Items      []AdvertisementResponse `json:"items"`
	HasNext    bool                    `json:"has_next"`
	NextCursor string                  `json:"next_cursor"`
	Total      int                     `json:"total"`
}
