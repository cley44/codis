package models

import (
	_ "ariga.io/atlas-provider-gorm/gormschema"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Username string
	Email    string
	Password string
	Test     string
}

func HashPassword(password string) (string, error) {
	return "", nil
}

// BeforeUpdate : hook before a user is updated
// func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
// 	if pw, err := bcrypt.GenerateFromPassword(u.Password, 0); err == nil {
// 		tx.Statement.SetColumn("Password", pw)
// 	}
// 	return
// }
