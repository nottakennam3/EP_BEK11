package types

import "time"

type UserSignupRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	UserProfile string `json:"userProfile"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Password    string `json:"password"`
	UserProfile string `json:"userProfile"`
}

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
	UserProfile string    `json:"userProfile"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewUser(username, password, profile string) *User {
	return &User{
		Username:    username,
		Password:    password,
		UserProfile: profile,
	}
}

type PostCreateRequest struct {
	Content string
}

type PostUpdateRequest struct {
	Content string
}

type PostCommentRequest struct {
	Content string
}

type Post struct {
	ID        int
	UserID    int
	Content   string
	CreatedAt time.Time
}

func NewPost(userID int, content string) *Post {
	return &Post{
		UserID: userID,
		Content: content,
	}
}

type PostComment struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	Timestamp time.Time
}

func NewPostComment(postID, userID int, content string) *PostComment {
	return &PostComment{
		PostID: postID,
		UserID: userID,
		Content: content,
	}
}

type PostLike struct {
	ID        int
	PostID    int
	UserID    int
	Timestamp time.Time
}
