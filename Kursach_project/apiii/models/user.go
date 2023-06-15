package models

type User struct {
	Userid   int64  `json:"userid,omitempty" db:"userid"`
	Login    string `json:"login,omitempty" db:"login"`
	Password string `json:"password,omitempty" db:"password"`
	Role     string `json:"role,omitempty" db:"role"`
}
