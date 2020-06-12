package handlers

import (
	"errors"
	"fmt"
	"hangout-api/context"
	"hangout-api/models"

	"github.com/gin-gonic/gin"
)

type Pub struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Image       string      `json:"image"`
	Geolocation Geolocation `json:"geolocation"`
	Detail      string      `json:"detail"`
}
type CreateReq struct {
	Name        string      `json:"name"`
	Image       string      `json:"image"`
	Detail      string      `json:"detail"`
	Geolocation Geolocation `json:"geolocation"`
}
type PubRes struct {
	Pub
}

func (h *Handler) CreatePub(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		c.Status(401)
		return
	}
	req := new(CreateReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	fmt.Println(req)
	pub := new(models.Pub)
	pub.Name = req.Name
	pub.Image = req.Image
	pub.Detail = req.Detail
	pub.UserID = user.ID
	err := h.ps.Create(pub)
	if err != nil {
		Error(c, 500, err)
		return
	}
	geo := new(models.Geolocation)
	geo.Longtitude = req.Geolocation.Longtitude
	geo.Latitude = req.Geolocation.Latitude
	geo.PubID = pub.ID
	err = h.gs.Create(geo)
	if err != nil {
		Error(c, 500, err)
		return
	}
	pubres, err := h.ps.GetByUserID(user.ID)
	if err != nil {
		Error(c, 500, err)
		return
	}
	res := new(PubRes)
	res.ID = pubres.ID
	res.Name = pubres.Name
	res.Image = pubres.Image
	res.Detail = pubres.Detail
	res.Geolocation.Longtitude = geo.Longtitude
	res.Geolocation.Latitude = geo.Latitude

	c.JSON(201, res)
}

func (h *Handler) ListPub(c *gin.Context) {
	data, err := h.ps.ListAllPub()
	if err != nil {
		Error(c, 500, err)
		return
	}
	pubs := []Pub{}
	for _, p := range data {
		geo := Geolocation{
			Longtitude: p.Geolocation.Longtitude,
			Latitude:   p.Geolocation.Latitude,
		}
		pubs = append(pubs, Pub{
			ID:          p.ID,
			Name:        p.Name,
			Image:       p.Image,
			Detail:      p.Detail,
			Geolocation: geo,
		})
	}
	c.JSON(200, pubs)
}

func (h *Handler) GetMyPub(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		Error(c, 401, errors.New("invalid token"))
		return
	}
	pub, err := h.ps.GetByUserID(user.ID)
	if err != nil {
		Error(c, 500, err)
		return
	}
	geo, err := h.gs.GetbyPubID(pub.ID)
	if err != nil {
		Error(c, 500, err)
		return
	}
	geores := Geolocation{
		Longtitude: geo.Longtitude,
		Latitude:   geo.Latitude,
	}
	pubres := &Pub{
		ID:          pub.ID,
		Name:        pub.Name,
		Image:       pub.Image,
		Detail:      pub.Detail,
		Geolocation: geores,
	}
	c.JSON(200, pubres)
}

type UpdateReq struct {
	Name        string      `json:"name"`
	Image       string      `json:"image"`
	Detail      string      `json:"detail"`
	Geolocation Geolocation `json:"geolocation"`
	PubID       uint        `json:"pub_id"`
}

func (h *Handler) UpdatePub(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		Error(c, 401, errors.New("invalid token"))
		return
	}
	req := new(UpdateReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	fmt.Println(req)
	pub := new(models.Pub)
	pub.Name = req.Name
	pub.Image = req.Image
	pub.Detail = req.Detail
	err := h.ps.UpdatePub(user.ID, pub)
	if err != nil {
		Error(c, 500, err)
		return
	}
	geo := &models.Geolocation{
		Longtitude: req.Geolocation.Longtitude,
		Latitude:   req.Geolocation.Latitude,
	}
	fmt.Println(req.PubID)
	err = h.gs.UpdateGeo(req.PubID, geo)
	if err != nil {
		Error(c, 500, err)
		return
	}
	geores := Geolocation{
		Longtitude: req.Geolocation.Longtitude,
		Latitude:   req.Geolocation.Latitude,
	}

	pubres := &Pub{
		ID:          req.PubID,
		Name:        pub.Name,
		Image:       pub.Image,
		Detail:      pub.Detail,
		Geolocation: geores,
	}
	c.JSON(200, pubres)
}
