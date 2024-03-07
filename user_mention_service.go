package main

import (
	"context"
	"fmt"

	"github.com/ServiceWeaver/weaver"
)

type IUserMentionService interface {
	ComposeUserMentions(context.Context, []string) []UserMention
}

type UserMentionService struct {
	weaver.Implements[IUserMentionService]

	storage weaver.Ref[Storage]
}

func (s *UserMentionService) ComposeUserMentions(ctx context.Context, usernames []string) []UserMention {
	storage := s.storage.Get()
	user_mentions := make([]UserMention, 0)
	for i, username := range usernames {
		user_profile, exist := storage.GetUserProfile(ctx, username)
		if !exist {
			fmt.Printf("[ComposeUserMentions] User profile not found for username: %s\n", username)
		} else {
			user_mentions = append(user_mentions, UserMention{
				userId:   user_profile.userId,
				username: usernames[i],
			})
		}
	}
	return user_mentions
}
