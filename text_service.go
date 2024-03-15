package main

import (
	"context"
	"regexp"

	"github.com/ServiceWeaver/weaver"
)

type ITextService interface {
	ComposeText(context.Context, string) TextServiceReturn
}

type TextService struct {
	weaver.Implements[ITextService]
	url_shorten_service    weaver.Ref[UrlShortenService]
	user_mention_service    weaver.Ref[UserMentionService]
}

func iterative_search(text string, pattern string) []string {
	res := make([]string, 0)
	reg := regexp.MustCompile(pattern)
	for i:=0 ; i < len(text) ; {
		loc := reg.FindIndex([]byte(text[i:]))
		st, ed := loc[0], loc[1]
		if st == -1 {
			break
		}
		res = append(res, text[i+st:i+ed])
		i += ed
	}
	return res
}


func (s *TextService) ComposeText(ctx context.Context, text string) TextServiceReturn {
	url_pattern := "(http://|https://)([a-zA-Z0-9_!~*'().&=+$%-]+)"
	mention_pattern :=  "@[a-zA-Z0-9-_]+"
	
	// regex search text for urls
	urls := iterative_search(text, url_pattern)

	// shorten urls
	url_shorten_service := s.url_shorten_service.Get()
	shorten_urls := make([]Url, 0)
	for _, url := range urls {
		shorten_urls = append(shorten_urls, url_shorten_service.ComposeUrl(ctx, []string{url})[0])
	}
	
	// regex search text for mentions
	mentions_str := iterative_search(text, mention_pattern)

	// convert mentions to UserMention type
	user_mention_service := s.user_mention_service.Get()
	mentions := user_mention_service.ComposeUserMentions(ctx, mentions_str)

	ret := TextServiceReturn{
		text: text,
		user_mentions: mentions,
		urls: shorten_urls,
	}

	if len(urls) > 0 {
		reg := regexp.MustCompile(url_pattern)
		updated_text := ""
		for i:=0 ; i < len(text) ; {
			loc := reg.FindIndex([]byte(text[i:]))
			st, ed := loc[0], loc[1]
			if st == -1 {
				updated_text += text[i:]
				break
			}
			updated_text += text[i:i+st] + shorten_urls[0].shortenedUrl
			i += ed
		}
		ret.text = updated_text
	}
	return ret
}