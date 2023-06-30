package user_dot

import "time"

type UserDOT struct {
	UserID      int32      `json:"userid" gorm:"column:userid;"`
	Email       string     `json:"email" gorm:"column:email;"`
	Password    string     `json:"password" gorm:"column:password;"`
	PhoneNumber int64      `json:"phone_number" gorm:"column:phone_number;"`
	FirstName   string     `json:"first_name" gorm:"column:first_name;"`
	LastName    string     `json:"last_name" gorm:"column:last_name;"`
	Verified    int8       `json:"verified" gorm:"column:verified;"`
	AvatarCDN   int32      `json:"avatar_cdn" gorm:"column:avatar_cdn;"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at;"`
}

type UserLogInDOT struct {
	Email    string `json:"email" gorm:"column:email;"`
	Password string `json:"password" gorm:"column:password;"`
}

type UserResponseDOT struct {
	UserID      int32  `json:"userid" gorm:"column:userid;"`
	Email       string `json:"email" gorm:"column:email;"`
	PhoneNumber int64  `json:"phone_number" gorm:"column:phone_number;"`
	FirstName   string `json:"first_name" gorm:"column:first_name;"`
	LastName    string `json:"last_name" gorm:"column:last_name;"`
	Verified    int8   `json:"verified" gorm:"column:verified;"`
}
