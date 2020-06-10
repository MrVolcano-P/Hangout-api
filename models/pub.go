package models

import "github.com/jinzhu/gorm"

type Pub struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Image       string `gorm:"not null"`
	Geolocation Geolocation
	Review      []Review
}
type PubService interface {
	Create(pub *Pub) error
	ListAllPub() ([]Pub, error)
	// Login(user *User) (string, error)
	// GetByToken(token string) (*User, error)
	// Logout(user *User) error
	// GetByID(id uint) (*User, error)
	// UpdateProfile(id uint, name string) error
	// CheckUsername(username string) bool
}

func NewPubService(db *gorm.DB) PubService {
	return &pubGorm{db}
}

type pubGorm struct {
	db *gorm.DB
}

func (pg *pubGorm) Create(pub *Pub) error {
	return pg.db.Create(pub).Error
}
func (pg *pubGorm) ListAllPub() ([]Pub, error) {
	pubs := []Pub{}
	err := pg.db.Find(&pubs).Error
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(pubs); i++ {
		geo := Geolocation{}
		err := pg.db.
			Where("pub_id = ?", pubs[i].ID).
			Find(&geo).Error
		if err != nil {
			return nil, err
		}
		pubs[i].Geolocation = geo
	}
	return pubs, nil
}
