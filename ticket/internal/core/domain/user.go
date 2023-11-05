package domain

import "time"

type User struct {
	BaseUser
	Password string    `json:"-"`
	Role     string    `json:"-"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

type BaseUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageUrl string `json:"image_url"`
}
