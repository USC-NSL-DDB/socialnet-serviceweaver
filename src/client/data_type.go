package main

import "github.com/ServiceWeaver/weaver"

type Creator struct {
	weaver.AutoMarshal
	UserId   int64
	Username string
}

type UserProfile struct {
	weaver.AutoMarshal
	UserId         int64
	FirstName      string
	LastName       string
	Salt           string
	PasswordHashed string
}

type PostType int

const (
	POST   PostType = 0
	REPOST PostType = 1
	REPLY  PostType = 2
	DM     PostType = 3
)

type Post struct {
	weaver.AutoMarshal
	Post_id       int64
	Creator       Creator
	Req_id        int64
	Text          string
	User_mentions []UserMention
	Media         []Media
	Urls          []Url
	Timestamp     int64
	Post_type     PostType
}

type Media struct {
	weaver.AutoMarshal
	MediaId   int64
	MediaType string
}

type Url struct {
	weaver.AutoMarshal
	ShortenedUrl string
	ExpandedUrl  string
}

type UserMention struct {
	weaver.AutoMarshal
	UserId   int64
	Username string
}

type TextServiceReturn struct {
	weaver.AutoMarshal
	Text          string
	User_mentions []UserMention
	Urls          []Url
}
