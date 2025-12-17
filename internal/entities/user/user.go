package user

import (
    "time"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Role         string    `json:"role"`
    Name         string    `json:"name"`
    Email        string    `gorm:"uniqueIndex" json:"email"`
    PasswordHash string    `json:"-"`
    Address      string    `json:"address"`
    Phone        string    `json:"phone"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) SetPassword(plain string) error {
    h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil { return err }
    u.PasswordHash = string(h)
    return nil
}

func (u *User) CheckPassword(plain string) bool {
    return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain)) == nil
}

