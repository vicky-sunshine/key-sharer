package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned when a resource cannot be found in the database.
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidPassword is returned when an invalid password // is used when attempting to authenticate a user.
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

// User define basic user info in database
type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Name         string
}

// UserService handle database connection with User resourse
type UserService struct {
	db            *gorm.DB
	servicePepper string
}

// NewUserService init the UserService database connection
func NewUserService(dbType string, connectionInfo string, servicePepper string) (*UserService, error) {
	db, err := gorm.Open(dbType, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db:            db,
		servicePepper: servicePepper,
	}, nil
}

// AutoMigrate will attempt to automatically migrate the
// users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// Close closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// ByUsername get User resourece by usereame
func (us *UserService) ByUsername(username string) (*User, error) {

	var user User
	db := us.db.Where("username=?", username)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create will create the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + us.servicePepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

// Update will update the provided user with all of the data
// in the provided user object.
func (us *UserService) Update(user *User) error {
	// find user by Username and get its ID
	var foundUser User
	db := us.db.Where("username=?", user.Username)
	err := first(db, &foundUser)
	if err != nil {
		return err
	}

	user.ID = foundUser.ID
	return us.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (us *UserService) Delete(username string) error {
	var user User
	db := us.db.Where("username=?", username)
	err := first(db, &user)
	if err != nil {
		return err
	}

	return us.db.Delete(&user).Error
}

// Authenticate can be used to authenticate a user with the
// provided username and password.
// If the username provided is invalid, this will return (nil, ErrNotFound)
// If the password provided is invalid, this will return (nil, ErrInvalidPassword)
// If the username and password are both valid, this will return (user, nil)
// Otherwise if another error is encountered this will return (nil, error)
func (us *UserService) Authenticate(username, password string) (*User, error) {
	foundUser, err := us.ByUsername(username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+us.servicePepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
}

// first will query using the provided gorm.DB and
// it will get the first item returned and place it into dst.
// If nothing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
