package handlers

import (
	"errors"
	"fmt"
	"hangout-api/context"
	"hangout-api/models"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	DOB       time.Time `json:"dob"`
}
type SignupReq struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	DOB       string `json:"dob"`
}
type UserParty struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

func (h *Handler) Signup(c *gin.Context) {
	req := new(SignupReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	user := new(models.User)
	user.Username = req.Username
	user.Email = req.Email
	user.Password = req.Password
	user.Name = req.Name
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	t1, e := time.Parse(
		time.RFC3339,
		req.DOB)
	if e != nil {
		Error(c, 500, e)
	}
	user.DOB = t1
	// fmt.Println(user.DOB)
	if err := h.us.Create(user); err != nil {
		Error(c, 500, err)
		return
	}
	c.JSON(201, gin.H{
		"token":     user.Token,
		"username":  user.Username,
		"email":     user.Email,
		"name":      user.Name,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"dob":       user.DOB,
	})
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *gin.Context) {
	req := new(LoginReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	user := new(models.User)
	user.Username = req.Username
	user.Password = req.Password
	token, err := h.us.Login(user)
	if err != nil {
		Error(c, 401, err)
		return
	}
	c.JSON(201, gin.H{
		"token": token,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		Error(c, 401, errors.New("invalid token"))
		return
	}
	err := h.us.Logout(user)
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(204)
}

func (h *Handler) GetProfile(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		Error(c, 401, errors.New("invalid token"))
		return
	}
	image := h.ims.GetByUserID(user.ID)
	if image == nil {
		c.JSON(200, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
			"fistName": user.FirstName,
			"lastName": user.LastName,
			"dob":      user.DOB,
			"image":    "",
		})
	} else {
		c.JSON(200, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
			"fistName": user.FirstName,
			"lastName": user.LastName,
			"dob":      user.DOB,
			"image":    path.Join(models.UploadPath, image.Filename),
		})
	}
}

type UpdateProfileReq struct {
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Name      string    `json:"name"`
	DOB       time.Time `json:"dob"`
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	user := context.User(c)
	if user == nil {
		Error(c, 401, errors.New("invalid token"))
		return
	}
	req := new(UpdateProfileReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}

	userReq := new(models.User)
	userReq.Name = req.Name
	userReq.Email = req.Email
	userReq.FirstName = req.FirstName
	userReq.LastName = req.LastName
	userReq.DOB = req.DOB

	err := h.us.UpdateProfile(user.ID, userReq)
	if err != nil {
		Error(c, 500, err)
		return
	}
	// userUp, err := h.us.GetByID(uint(user.ID))
	// if err != nil {
	// 	Error(c, 500, err)
	// 	return
	// }
	c.JSON(200, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"name":     userReq.Name,
		"email":    userReq.Email,
		"fistName": userReq.FirstName,
		"lastName": userReq.LastName,
		"dob":      userReq.DOB,
		"image":    user.Image,
	})
}

func (h *Handler) CheckUsername(c *gin.Context) {
	username := c.Param("username")
	fmt.Println(username)
	status := h.us.CheckUsername(username)
	fmt.Println(status)
	c.JSON(200, gin.H{
		"is_Available": status,
	})
}
