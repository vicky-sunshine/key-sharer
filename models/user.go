package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// User define basic user info in database
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
}

var (
	// ErrNotFound is returned when a resource cannot be found in the database.
	ErrNotFound = errors.New("models: resource not found")
)

// UserService handle database connection with User resourse
type UserService struct {
	db *gorm.DB
}

// AutoMigrate will attempt to automatically migrate the
// users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// NewUserService init the UserService database connection
func NewUserService(dbType string, connectionInfo string) (*UserService, error) {
	db, err := gorm.Open(dbType, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil
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
