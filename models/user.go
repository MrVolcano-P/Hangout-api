package models

import (
	"fmt"
	"hangout-api/hash"
	"hangout-api/rand"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const cost = 12

type User struct {
	gorm.Model
	Username  string    `gorm:"unique_index;not null"`
	Email     string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Name      string    `gorm:"not null"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	DOB       time.Time `gorm:"not null"`
	Image     string
	Token     string `gorm:"index"`
	Role      string `gorm:"not null"`
	Review    []Review
}

type UserService interface {
	Create(user *User) error
	Login(user *User) (string, error)
	GetByToken(token string) (*User, error)
	Logout(user *User) error
	GetByID(id uint) (*User, error)
	UpdateProfile(id uint, user *User) error
	CheckUsername(username string) bool
	UpdateProfileImage(id uint, image string) error
}

func NewUserService(db *gorm.DB, hmac *hash.HMAC) UserService {
	return &userGorm{db, hmac}
}

type userGorm struct {
	db   *gorm.DB
	hmac *hash.HMAC
}

func (ug *userGorm) Create(temp *User) error {
	user := new(User)
	user.Email = temp.Email
	user.Password = temp.Password
	user.Name = temp.Name
	user.Username = temp.Username
	user.FirstName = temp.FirstName
	user.LastName = temp.LastName
	user.DOB = temp.DOB
	user.Role = temp.Role
	user.Image = temp.Image
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	token, err := rand.GetToken()
	if err != nil {
		return err
	}

	fmt.Println("token ===> ", token)
	tokenHash := ug.hmac.Hash(token)
	fmt.Println("tokenHashStr ===> ", tokenHash)

	user.Token = tokenHash
	temp.Token = token

	return ug.db.Create(user).Error
}

func (ug *userGorm) Login(user *User) (string, error) {
	found := new(User)
	err := ug.db.Where("username = ?", user.Username).First(&found).Error
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token, err := rand.GetToken()
	if err != nil {
		return "", err
	}

	fmt.Println("token ===> ", token)
	tokenHash := ug.hmac.Hash(token)
	fmt.Println("tokenHashStr ===> ", tokenHash)

	err = ug.db.Model(&User{}).
		Where("id = ?", found.ID).
		Update("token", tokenHash).Error
	if err != nil {
		return "", err
	}
	return token, nil
}

func (ug *userGorm) Logout(user *User) error {
	return ug.db.Model(user).
		Where("id = ?", user.ID).
		Update("token", "").Error
}

func (ug *userGorm) GetByToken(token string) (*User, error) {
	tokenHash := ug.hmac.Hash(token)
	log.Println("lookup for user by token(hashed): ", tokenHash)
	user := new(User)
	err := ug.db.Where("token = ?", tokenHash).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ug *userGorm) GetByID(id uint) (*User, error) {
	user := new(User)
	if err := ug.db.First(user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ug *userGorm) UpdateProfile(id uint, user *User) error {
	return ug.db.Model(&User{}).Where("id = ?", id).
		Updates(map[string]interface{}{"name": user.Name, "email": user.Email, "first_name": user.FirstName, "last_name": user.LastName, "dob": user.DOB}).Error
}
func (ug *userGorm) UpdateProfileImage(id uint, image string) error {
	return ug.db.Model(&User{}).Where("id = ?", id).
		Update("image", image).Error
}
func (ug *userGorm) CheckUsername(username string) bool {
	user := new(User)
	if err := ug.db.Where("username = ?", username).First(user).Error; gorm.IsRecordNotFoundError(err) {
		// record not found
		return true
	}
	return false
}
