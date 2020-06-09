package models

import "github.com/jinzhu/gorm"

type Geolocation struct {
	Longtitude string `gorm:"not null"`
	Latitude   string `gorm:"not null"`
	PubID      uint   `gorm:"not null"`
}
type GeoService interface {
	Create(geolocation *Geolocation) error
	// Login(user *User) (string, error)
	// GetByToken(token string) (*User, error)
	// Logout(user *User) error
	// GetByID(id uint) (*User, error)
	// UpdateProfile(id uint, name string) error
	// CheckUsername(username string) bool
}

func NewGeoService(db *gorm.DB) GeoService {
	return &geoGorm{db}
}

type geoGorm struct {
	db *gorm.DB
}

func (gg *geoGorm) Create(geolocation *Geolocation) error {
	return gg.db.Create(geolocation).Error
}
