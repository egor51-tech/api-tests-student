package models

import (
	advModels "api-tests-template/internal/managers/advertisements/models"
)

// FavoriteResponse запись избранного: приходит при добавлении и как элемент списка
type FavoriteResponse struct {
	FavoriteId    string                          `json:"favorite_id"`
	AddedAt       string                          `json:"added_at"`
	Advertisement advModels.AdvertisementResponse `json:"advertisement"`
}

type FavoriteListResponse struct {
	Items      []FavoriteResponse `json:"items"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	Total      int                `json:"total"`
	TotalPages int                `json:"total_pages"`
}

type BulkAddFavoritesResponse struct {
	Added        []string `json:"added"`
	Skipped      []string `json:"skipped"`
	NotFound     []string `json:"not_found"`
	LimitReached bool     `json:"limit_reached"`
}
