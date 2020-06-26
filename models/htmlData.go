package models

import (
	"html/template"
	)

type PostPage struct {
	UsernamePass string
	ErrorMsgPass template.HTML
}

type YourPostPage struct {
	UsernamePass 	string
	PostsWithEncode	[]PostWithEncode
}

type FindPostPage struct {
	UsernamePass 	string
	UserToFind		string
	PostsWithEncode	[]PostWithEncode
}