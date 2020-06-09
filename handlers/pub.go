package handlers

import (
	"fmt"
	"hangout-api/models"

	"github.com/gin-gonic/gin"
)

type Pub struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Image       string      `json:"image"`
	Geolocation Geolocation `json:"geolocation"`
}
type CreateReq struct {
	Name        string      `json:"name"`
	Image       string      `json:"image"`
	Geolocation Geolocation `json:"geolocation"`
}
type PubRes struct {
	Pub
}

func (h *Handler) CreatePub(c *gin.Context) {
	req := new(CreateReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	fmt.Println(req)
	pub := new(models.Pub)
	pub.Name = req.Name
	pub.Image = req.Image
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

	res := new(PubRes)
	res.ID = pub.ID
	res.Name = pub.Name
	res.Image = pub.Image
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
			Geolocation: geo,
		})
	}
	c.JSON(200, pubs)
}
