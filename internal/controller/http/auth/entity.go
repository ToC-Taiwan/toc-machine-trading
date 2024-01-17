package auth

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseBody struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
	Code   int    `json:"code"`
}
