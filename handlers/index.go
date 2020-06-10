package handlers

import "hangout-api/models"

type Handler struct {
	us  models.UserService
	ims models.ImageService
	ps  models.PubService
	gs  models.GeoService
	rs  models.ReviewService
}

func NewHandler(us models.UserService, ims models.ImageService, ps models.PubService, gs models.GeoService, rs models.ReviewService) *Handler {
	return &Handler{us, ims, ps, gs, rs}
}
