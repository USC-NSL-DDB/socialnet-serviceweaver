package main

import "github.com/ServiceWeaver/weaver"

type Creator struct {
	weaver.AutoMarshal
	userId   int64
	username string
}

type UserProfile struct {
	weaver.AutoMarshal
	userId         int64
	firstName      string
	lastName       string
	salt           string
	passwordHashed string
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
	post_id       int64
	creator       Creator
	req_id        int64
	text          string
	user_mentions []UserMention
	media         []Media
	urls          []Url
	timestamp     int64
	post_type     PostType
}

type Media struct {
	weaver.AutoMarshal
	mediaId   int64
	mediaType string
}

type Url struct {
	weaver.AutoMarshal
	shortenedUrl string
	expandedUrl  string
}

type UserMention struct {
	weaver.AutoMarshal
	userId   int64
	username string
}

type TextServiceReturn struct {
	weaver.AutoMarshal
	text          string
	user_mentions []UserMention
	urls          []Url
}
