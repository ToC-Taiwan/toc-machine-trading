package entity

type User struct {
	ID            int    `json:"-"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	PushToken     string `json:"push_token"`
	EmailVerified bool   `json:"-"`
	AuthTrade     bool   `json:"-"`
}
