package models

import (
	"github.com/jinzhu/gorm"
)

type Review struct {
	gorm.Model
	Text string `gorm:"not null"`
	// Date   time.Time `gorm:"not null"`
	UserID uint `gorm:"not null"`
	PubID  uint `gorm:"not null"`
}

type ReviewService interface {
	Create(review *Review) error
	GetByPubID(id uint) ([]Review, error)
	// ListAllPub() ([]Pub, error)
	// Login(user *User) (string, error)
	// GetByToken(token string) (*User, error)
	// Logout(user *User) error
	// GetByID(id uint) (*User, error)
	// UpdateProfile(id uint, name string) error
	// CheckUsername(username string) bool
}

func NewReviewService(db *gorm.DB) ReviewService {
	return &reviewGorm{db}
}

type reviewGorm struct {
	db *gorm.DB
}

func (rg *reviewGorm) Create(review *Review) error {
	return rg.db.Create(review).Error
}

func (rg *reviewGorm) GetByPubID(id uint) ([]Review, error) {
	reviews := []Review{}
	err := rg.db.
		Where("pub_id = ?", id).
		Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
