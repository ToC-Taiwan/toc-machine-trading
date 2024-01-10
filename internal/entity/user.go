package entity

type User struct {
	ID            int    `json:"-"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"-"`
	AuthTrade     bool   `json:"-"`
}
