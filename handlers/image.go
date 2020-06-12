package handlers

import (
	"fmt"
	"hangout-api/context"
	"hangout-api/models"
	"path"

	"github.com/gin-gonic/gin"
)

type ImageRes struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Filename string `json:"filename"`
}

type CreateImageRes struct {
	ImageRes
}

func (h *Handler) CreateImage(c *gin.Context) {
	userContext := context.User(c)
	if userContext == nil {
		c.Status(401)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		Error(c, 400, err)
		return
	}

	images, err := h.ims.CreateImages(form.File["photos"], userContext.ID)
	if err != nil {
		Error(c, 500, err)
		return
	}

	res := []CreateImageRes{}
	for _, img := range images {
		r := CreateImageRes{}
		r.ID = img.ID
		r.UserID = userContext.ID
		r.Filename = path.Join(models.UploadPath, img.Filename)
		res = append(res, r)
	}

	c.JSON(201, res)
}

type ListGalleryImagesRes struct {
	ImageRes
}

func (h *Handler) UpdateProfileImageTable(c *gin.Context) {
	userContext := context.User(c)
	if userContext == nil {
		c.Status(401)
		return
	}
	fmt.Println("in")
	// fmt.Println(userContext)
	form, err := c.MultipartForm()
	if err != nil {
		Error(c, 400, err)
		return
	}
	image := h.ims.GetByUserID(userContext.ID)
	fmt.Println(form.File["photos"])
	if image == nil {
		// Error(c, 500, err)
		// return
		images, err := h.ims.CreateImages(form.File["photos"], userContext.ID)
		if err != nil {
			Error(c, 500, err)
			return
		}
		// fmt.Println(images[0])
		r := CreateImageRes{}
		r.ID = images[0].ID
		r.UserID = userContext.ID
		r.Filename = path.Join(models.UploadPath, images[0].Filename)
		c.JSON(201, r)
	} else {
		err := h.ims.RemoveImageByFileName(image.Filename)
		if err != nil {
			Error(c, 500, err)
			return
		}
		images, err := h.ims.UpdateImage(form.File["photos"], userContext.ID)
		if err != nil {
			Error(c, 500, err)
			return
		}
		err = h.ims.UpdateProfileImg(userContext.ID, images[0].Filename)
		if err != nil {
			Error(c, 500, err)
			return
		}
		c.Status(204)
	}
}
