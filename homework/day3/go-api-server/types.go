package main

type CreateUserReq struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
	UserProfile	string	`json:"userProfile"`
}

type LoginReq struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}

type UpdateReq struct {
	Password	string	`json:"password"`
	UserProfile	string	`json:"userProfile"`
}

type User struct {
	ID			int		`json:"id"`
	Username	string	`json:"username"`
	EncPassword	string	`json:"password"`
	UserProfile	string	`json:"userProfile"`
}

func NewUser(username, password, userProfile string) *User {
	return &User{
		Username: username,
		EncPassword: password,
		UserProfile: userProfile,
	}
}
