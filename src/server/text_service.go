package main

import (
	"context"
	"regexp"

	"github.com/ServiceWeaver/weaver"
)

type ITextService interface {
	ComposeText(context.Context, string) (TextServiceReturn, error)
}

type TextService struct {
	weaver.Implements[ITextService]
	url_shorten_service  weaver.Ref[IUrlShortenService]
	user_mention_service weaver.Ref[IUserMentionService]
}

func iterative_search(text string, pattern string) []string {
	res := make([]string, 0)
	reg := regexp.MustCompile(pattern)
	for i := 0; i < len(text); {
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

func (s *TextService) ComposeText(ctx context.Context, text string) (TextServiceReturn, error) {
	url_pattern := "(http://|https://)([a-zA-Z0-9_!~*'().&=+$%-]+)"
	mention_pattern := "@[a-zA-Z0-9-_]+"

	// regex search text for urls
	urls := iterative_search(text, url_pattern)

	// shorten urls
	url_shorten_service := s.url_shorten_service.Get()
	new_urls, _ := url_shorten_service.ComposeUrl(ctx, urls)

	// regex search text for mentions
	mentions_str := iterative_search(text, mention_pattern)

	// convert mentions to UserMention type
	user_mention_service := s.user_mention_service.Get()
	mentions, _ := user_mention_service.ComposeUserMentions(ctx, mentions_str)

	ret := TextServiceReturn{
		Text:          text,
		User_mentions: mentions,
		Urls:          new_urls,
	}

	if len(urls) > 0 {
		reg := regexp.MustCompile(url_pattern)
		updated_text := ""
		for i := 0; i < len(text); {
			loc := reg.FindIndex([]byte(text[i:]))
			st, ed := loc[0], loc[1]
			if st == -1 {
				updated_text += text[i:]
				break
			}
			updated_text += text[i:i+st] + new_urls[0].ShortenedUrl
			i += ed
		}
		ret.Text = updated_text
	}
	return ret, nil
}
