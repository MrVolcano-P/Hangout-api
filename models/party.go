package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Party struct {
	gorm.Model
	Name       string    `gorm:"not null"`
	Membership string    `gorm:"not null"`
	Date       time.Time `gorm:"not null"`
	PubID      uint      `gorm:"not null`
	UserID     uint      `gorm:"not null"`
	Members    []Member
}
type PartyService interface {
	Create(party *Party) error
	GetPartiesByPubId(id uint) ([]Party, error)
	GetPartiesByUserId(id uint) ([]Party, error)
	GetPartyById(id uint) (*Party, error)
}

func NewPartyService(db *gorm.DB) PartyService {
	return &partyGorm{db}
}

type partyGorm struct {
	db *gorm.DB
}

func (pg *partyGorm) Create(party *Party) error {
	return pg.db.Create(party).Error
}

func (pg *partyGorm) GetPartiesByPubId(id uint) ([]Party, error) {
	parties := []Party{}
	err := pg.db.Where("pub_id = ?", id).Find(&parties).Error
	if err != nil {
		return nil, err
	}
	return parties, nil
}

func (pg *partyGorm) GetPartiesByUserId(id uint) ([]Party, error) {
	parties := []Party{}
	err := pg.db.Where("user_id = ?", id).Find(&parties).Error
	if err != nil {
		return nil, err
	}
	return parties, nil
}

func (pg *partyGorm) GetPartyById(id uint) (*Party, error) {
	party := new(Party)
	err := pg.db.Where("id = ?", id).First(party).Error
	if err != nil {
		return nil, err
	}
	return party, nil
}
