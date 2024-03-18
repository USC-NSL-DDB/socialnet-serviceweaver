package main

import (
	"context"
	"fmt"

	"github.com/ServiceWeaver/weaver"
)

type IUserMentionService interface {
	ComposeUserMentions(context.Context, []string) ([]UserMention, error)
}

type UserMentionService struct {
	weaver.Implements[IUserMentionService]

	storage weaver.Ref[IStorage]
}

func (s *UserMentionService) ComposeUserMentions(ctx context.Context, usernames []string) ([]UserMention, error) {
	storage := s.storage.Get()
	user_mentions := make([]UserMention, 0)
	for i, username := range usernames {
		user_profile, exist, _ := storage.GetUserProfile(ctx, username)
		if !exist {
			fmt.Printf("[ComposeUserMentions] User profile not found for username: %s\n", username)
		} else {
			user_mentions = append(user_mentions, UserMention{
				UserId:   user_profile.UserId,
				Username: usernames[i],
			})
		}
	}
	return user_mentions, nil
}
