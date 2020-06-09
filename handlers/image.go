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

// type Handler struct {
// 	gs  models.GalleryService
// 	ims models.ImageService
// }

// func NewHandler(gs models.GalleryService, ims models.ImageService) *Handler {
// 	return &Handler{gs, ims}
// }

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

// func (h *Handler) ListGalleryImages(c *gin.Context) {
// 	galleryIDStr := c.Param("id")
// 	id, err := strconv.Atoi(galleryIDStr)
// 	if err != nil {
// 		Error(c, 400, err)
// 		return
// 	}

// 	gallery, err := h.gs.GetByID(uint(id))
// 	if err != nil {
// 		Error(c, 400, err)
// 		return
// 	}
// 	images, err := h.ims.GetByGalleryID(gallery.ID)
// 	if err != nil {
// 		Error(c, http.StatusNotFound, err)
// 		return
// 	}
// 	res := []ListGalleryImagesRes{}
// 	for _, img := range images {
// 		r := ListGalleryImagesRes{}
// 		r.ID = img.ID
// 		r.GalleryID = gallery.ID
// 		r.Filename = img.FilePath()
// 		res = append(res, r)
// 	}
// 	c.JSON(http.StatusOK, res)
// }

// type DeleteReq struct {
// 	FileNames []string `json:"filenames"`
// }

// func (h *Handler) DeleteImageInGallary(c *gin.Context) {
// 	galleryIDStr := c.Param("id")
// 	id, err := strconv.Atoi(galleryIDStr)
// 	if err != nil {
// 		Error(c, 400, err)
// 		return
// 	}
// 	req := new(DeleteReq)
// 	if err := c.BindJSON(req); err != nil {
// 		Error(c, 400, err)
// 		return
// 	}
// 	fmt.Println(req)
// 	for _, r := range req.FileNames {
// 		fmt.Println(r)
// 		err = h.ims.RemoveImageByFileName(r)
// 		if err != nil {
// 			Error(c, 500, err)
// 			return
// 		}
// 	}
// 	c.Status(204)
// }

// func (h *Handler) DeleteImage(c *gin.Context) {
// 	imageIDStr := c.Param("id")
// 	filename := c.Param("filename")
// 	id, err := strconv.Atoi(imageIDStr)
// 	if err != nil {
// 		Error(c, 400, err)
// 		return
// 	}
// 	err = h.ims.RemoveImageByFileName(filename)
// 	if err != nil {
// 		Error(c, 500, err)
// 		return
// 	}
// 	c.Status(200)
// 	// if err := h.ims.Delete(uint(id)); err != nil {
// 	// 	Error(c, 500, err)
// 	// 	return
// 	// }
// }
func (h *Handler) UpdateProfileImage(c *gin.Context) {
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
