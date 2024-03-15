package main

type Creator struct {
	userId   int64
	username string
}

type UserProfile struct {
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
	post_id int64
	creator Creator
	req_id  int64
	text    string
	// user_mentions []UserMention
	// media []Media
	// urls []URL
	timestamp int64
	post_type PostType
}

type Media struct {
	mediaId   int64
	mediaType string
}

type Url struct {
	shortenedUrl string
	expandedUrl  string
}

type UserMention struct {
	userId  int64
	username string
}

type TextServiceReturn struct {
	text string
	user_mentions []UserMention
	urls []Url
}