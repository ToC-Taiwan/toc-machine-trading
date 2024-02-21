package entity

import "time"

type AuthJWT struct {
	ID      int
	Key     string
	Created time.Time
}
