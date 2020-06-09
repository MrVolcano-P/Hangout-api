package main

import (
	"hangout-api/config"
	"hangout-api/handlers"
	"hangout-api/hash"
	"hangout-api/middleware"
	"hangout-api/models"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	conf := config.Load()

	db, err := gorm.Open("mysql", conf.Connection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if conf.Mode == "dev" {
		db.LogMode(true) // dev only!
	}

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
	}

	hmac := hash.NewHMAC(conf.HMACKey)
	us := models.NewUserService(db, hmac)
	ims := models.NewImageService(db)

	h := handlers.NewHandler(us, ims)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if conf.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Static("/upload", "./upload")
	r.POST("/signup", h.Signup)
	r.POST("/login", h.Login)
	r.GET("/checkUsername/:username", h.CheckUsername)
	auth := r.Group("/")
	auth.Use(middleware.RequireUser(us))
	{
		auth.POST("/logout", h.Logout)
		user := auth.Group("/user")
		{
			user.GET("/profile", h.GetProfile)
			user.POST("/profile", h.UpdateProfile)
			user.POST("/profile/image", h.CreateImage)
			user.PUT("/profile/image", h.UpdateProfileImage)
		}
	}
	r.Run(":8080")
}
