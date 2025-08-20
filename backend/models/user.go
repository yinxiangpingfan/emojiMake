package models

import "golang.org/x/crypto/bcrypt"

// User 用户模型
type User struct {
	ID       int64  `xorm:"id pk autoincr" json:"id"`
	Phone    string `xorm:"phone unique" json:"phone"`
	Password string `xorm:"password" json:"-"` // 密码在json中不可见
}

// SetPassword 对密码进行加密
func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
