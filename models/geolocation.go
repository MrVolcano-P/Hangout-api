package models

import "github.com/jinzhu/gorm"

type Geolocation struct {
	Longtitude string `gorm:"not null"`
	Latitude   string `gorm:"not null"`
	PubID      uint   `gorm:"not null"`
}
type GeoService interface {
	Create(geolocation *Geolocation) error
	GetbyPubID(id uint) (*Geolocation, error)
	UpdateGeo(id uint, geo *Geolocation) error
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
func (gg *geoGorm) GetbyPubID(id uint) (*Geolocation, error) {
	geo := &Geolocation{}
	err := gg.db.Where("pub_id = ?", id).First(geo).Error
	if err != nil {
		return nil, err
	}
	return geo, err
}
func (gg *geoGorm) UpdateGeo(id uint, geo *Geolocation) error {
	return gg.db.Model(&Geolocation{}).Where("pub_id = ?", id).
		Updates(map[string]interface{}{"longtitude": geo.Longtitude, "latitude": geo.Latitude}).Error
}
