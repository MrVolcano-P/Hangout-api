package handlers

import "hangout-api/models"

type Handler struct {
	us models.UserService
	ims models.ImageService
}

func NewHandler(us models.UserService,ims models.ImageService) *Handler {
	return &Handler{us,ims}
}
