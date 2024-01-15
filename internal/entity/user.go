package entity

import "time"

type User struct {
	ID            int    `json:"-"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"-"`
	AuthTrade     bool   `json:"-"`
}

type PushToken struct {
	ID      int
	Token   string
	UserID  int
	Enabled bool
	Created time.Time
}
