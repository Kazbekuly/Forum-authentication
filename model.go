package Forum

import (
	"html/template"
	"time"
)

type User struct {
	Id       int
	Email    string
	Username string
	Password string
}

type Token struct {
	Id        int
	UserId    int
	AuthToken string
	ExpiresAT time.Time
}

type Post struct {
	Id         int
	UserId     int
	Title      string
	Text       string
	Categories string
	CreatedAt  string
	Author     string
	Like       string
	Dislike    string
	Image      template.URL
}

type MyPost struct {
	User  string
	Posts []Post
}

type UserCategory struct {
	Username string
	Category []Post
}

type Comments struct {
	Id        int
	UserId    int
	PostId    int
	Comment   string
	CreatedAt string
	Author    string
	Like      string
	Dislike   string
}

type GetComments struct {
	User     string
	Post     Post
	Comments []Comments
	Images   []template.URL
}

type Evaluate struct {
	Id     int
	PostId int
	UserId int
	Vote   int
}
type Vote struct {
	Like    string
	Dislike string
}

type EvaluateComment struct {
	Id        int
	UserId    int
	CommentId int
	Vote      int
}
