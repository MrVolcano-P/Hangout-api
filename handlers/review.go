package handlers

import (
	"hangout-api/context"
	"hangout-api/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Review struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
	// Date time.Time `json:"date"`
}
type ReviewReq struct {
	Text string `json:"text"`
}
type ReviewRes struct {
	ID   uint      `json:"id"`
	Text string    `json:"text"`
	Date time.Time `json:"date"`
	User UserRes   `json:"user"`
}
type UserRes struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

func (h *Handler) CreateReview(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		c.Status(401)
		return
	}
	pubId := c.Param("id")
	id, err := strconv.Atoi(pubId)
	if err != nil {
		Error(c, 400, err)
		return
	}
	req := new(Review)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	review := new(models.Review)
	review.Text = req.Text
	review.PubID = uint(id)
	review.UserID = user.ID

	err = h.rs.Create(review)
	if err != nil {
		Error(c, 500, err)
		return
	}
	userRes := UserRes{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
	}
	res := ReviewRes{
		ID:   review.ID,
		Text: review.Text,
		Date: review.CreatedAt,
		User: userRes,
	}
	c.JSON(200, res)
}

func (h *Handler) GetReviewByPubID(c *gin.Context) {
	pubId := c.Param("id")
	id, err := strconv.Atoi(pubId)
	if err != nil {
		Error(c, 400, err)
		return
	}
	reviews, err := h.rs.GetByPubID(uint(id))

	if err != nil {
		Error(c, 500, err)
		return
	}
	ress := []ReviewRes{}
	for _, review := range reviews {
		user, err := h.us.GetByID(review.UserID)
		if err != nil {
			Error(c, 500, err)
			return
		}
		userRes := UserRes{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
		}
		res := ReviewRes{
			ID:   review.ID,
			Text: review.Text,
			Date: review.CreatedAt,
			User: userRes,
		}
		ress = append(ress, res)
	}
	c.JSON(200, ress)
}
