package main
import (
	"context"
	"github.com/ServiceWeaver/weaver"
	"fmt"
)


type ISocialGraphService interface {
	GetFollowers(context.Context, int64) []int64
	GetFollowees(context.Context, int64) []int64
	Follow(context.Context, int64, int64)
	Unfollow(context.Context, int64, int64)
	FollowWithUsername(context.Context, string, string)
	UnfollowWithUsername(context.Context, string, string)
}

type SocialGraphService struct {
	weaver.Implements[ISocialGraphService]
	storage weaver.Ref[Storage]
	user_service weaver.Ref[UserService]
}

func map_to_list(m map[int64]bool) []int64 {
	l := make([]int64, 0, len(m))
	for k := range m {
		l = append(l, k)
	}
	return l
}


func (s *SocialGraphService) GetFollowers(ctx context.Context, userId int64) []int64 {
	storage := s.storage.Get()
	follower_maps, ok := storage.GetFollowers(ctx, userId)
	if !ok {
		return []int64{}
	}
	return map_to_list(follower_maps)
}

func (s *SocialGraphService) GetFollowees(ctx context.Context, userId int64) []int64 {
	storage := s.storage.Get()
	followee_maps, ok := storage.GetFollowees(ctx, userId)
	if !ok {
		return []int64{}
	}
	return map_to_list(followee_maps)
}

func (s *SocialGraphService) Follow(ctx context.Context, followerId int64, followeeId int64) {
	storage := s.storage.Get()
	storage.Follow(ctx, followerId, followeeId)
}

func (s *SocialGraphService) Unfollow(ctx context.Context, followerId int64, followeeId int64) {
	storage := s.storage.Get()
	storage.Unfollow(ctx, followerId, followeeId)
}

func (s *SocialGraphService) FollowWithUsername(ctx context.Context, followerUsername string, followeeUsername string) {
	user_service := s.user_service.Get()
	followerId := user_service.GetUserId(ctx, followerUsername) 
	followeeId := user_service.GetUserId(ctx, followeeUsername)
	if followerId == 0 || followeeId == 0 {
		fmt.Printf("Failed to find the user profile - followerUsername: %s, followeeUsername: %s\n", followerUsername, followeeUsername)
		return
	}
	storage := s.storage.Get()
	storage.Follow(ctx, followerId, followeeId)
}

func (s *SocialGraphService) UnfollowWithUsername(ctx context.Context, followerUsername string, followeeUsername string) {
	user_service := s.user_service.Get()
	followerId := user_service.GetUserId(ctx, followerUsername)  // TODO: change to multithread fetching
	followeeId := user_service.GetUserId(ctx, followeeUsername)
	if followerId == 0 || followeeId == 0 {
		fmt.Printf("Failed to find the user profile - followerUsername: %s, followeeUsername: %s\n", followerUsername, followeeUsername)
		return
	}
	storage := s.storage.Get()
	storage.Unfollow(ctx, followerId, followeeId)
}