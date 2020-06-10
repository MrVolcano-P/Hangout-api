package handlers

import "hangout-api/models"

type Handler struct {
	us  models.UserService
	ims models.ImageService
	ps  models.PubService
	gs  models.GeoService
	rs  models.ReviewService
	pts models.PartyService
	ms  models.MemberService
}

func NewHandler(us models.UserService, ims models.ImageService, ps models.PubService, gs models.GeoService, rs models.ReviewService,
	pts models.PartyService, ms models.MemberService) *Handler {
	return &Handler{us, ims, ps, gs, rs, pts, ms}
}
