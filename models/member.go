package models

import "github.com/jinzhu/gorm"

type Member struct {
	PartyID uint `gorm:not null"`
	UserID  uint `gorm:"not null"`
}

type MemberService interface {
	Create(member *Member) error
	GetByPartyID(id uint) ([]Member, error)
	GetByUserID(id uint) ([]Member, error)
}

func NewMemberService(db *gorm.DB) MemberService {
	return &memberGorm{db}
}

type memberGorm struct {
	db *gorm.DB
}

func (mg *memberGorm) Create(member *Member) error {
	return mg.db.Create(member).Error
}
func (mg *memberGorm) GetByPartyID(id uint) ([]Member, error) {
	members := []Member{}
	err := mg.db.Where("party_id = ?", id).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (mg *memberGorm) GetByUserID(id uint) ([]Member, error) {
	members := []Member{}
	err := mg.db.Where("user_id = ?", id).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
