package handlers

import "hangout-api/models"

type Handler struct {
	us  models.UserService
	ims models.ImageService
	ps  models.PubService
	gs  models.GeoService
}

func NewHandler(us models.UserService, ims models.ImageService, ps models.PubService,gs models.GeoService) *Handler {
	return &Handler{us, ims, ps,gs}
}
