package entity

import "time"

type User struct {
	ID            int    `json:"-"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Password      string `json:"-"`
	EmailVerified bool   `json:"-"`
	AuthTrade     bool   `json:"-"`
}

type NewUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type PushToken struct {
	ID      int
	Token   string
	UserID  int
	Enabled bool
	Created time.Time
}
