package main

import (
	"context"
	"fmt"

	"SocialNetwork/shared/common"

	"github.com/ServiceWeaver/weaver"
)

type ISocialGraphService interface {
	GetFollowers(context.Context, int64) ([]int64, error)
	GetFollowees(context.Context, int64) ([]int64, error)
	Follow(context.Context, int64, int64) error
	Unfollow(context.Context, int64, int64) error
	FollowWithUsername(context.Context, string, string) error
	UnfollowWithUsername(context.Context, string, string) error
}

type SocialGraphService struct {
	weaver.Implements[ISocialGraphService]
	storage      weaver.Ref[IStorage]
	user_service weaver.Ref[UserServicer]
}

func map_to_list(m map[int64]bool) []int64 {
	l := make([]int64, 0, len(m))
	for k := range m {
		l = append(l, k)
	}
	return l
}

func (s *SocialGraphService) GetFollowers(ctx context.Context, userId int64) ([]int64, error) {
	storage := s.storage.Get()
	follower_maps, ok, _ := storage.GetFollowers(ctx, userId)
	if !ok {
		return []int64{}, nil
	}
	return map_to_list(follower_maps), nil
}

func (s *SocialGraphService) GetFollowees(ctx context.Context, userId int64) ([]int64, error) {
	storage := s.storage.Get()
	followee_maps, ok, _ := storage.GetFollowees(ctx, userId)
	if !ok {
		return []int64{}, nil
	}
	return map_to_list(followee_maps), nil
}

func (s *SocialGraphService) Follow(ctx context.Context, followerId int64, followeeId int64) error {
	storage := s.storage.Get()
	storage.Follow(ctx, followerId, followeeId)
	return nil
}

func (s *SocialGraphService) Unfollow(ctx context.Context, followerId int64, followeeId int64) error {
	storage := s.storage.Get()
	storage.Unfollow(ctx, followerId, followeeId)
	return nil
}

func (s *SocialGraphService) FollowWithUsername(ctx context.Context, followerUsername string, followeeUsername string) error {
	user_service := s.user_service.Get()
	followerId, _ := user_service.GetUserId(ctx, followerUsername)
	followeeId, _ := user_service.GetUserId(ctx, followeeUsername)
	if followerId == 0 || followeeId == 0 {
		fmt.Printf("Failed to find the user profile - followerUsername: %s, followeeUsername: %s\n", followerUsername, followeeUsername)
		return nil
	}
	storage := s.storage.Get()
	storage.Follow(ctx, followerId, followeeId)
	return nil
}

func (s *SocialGraphService) UnfollowWithUsername(ctx context.Context, followerUsername string, followeeUsername string) error {
	user_service := s.user_service.Get()
	follower_id_fu := common.AsyncExec(func() interface{} {
		followerId, _ := user_service.GetUserId(ctx, followerUsername)
		return followerId
	})

	followee_id_fu := common.AsyncExec(func() interface{} {
		followeeId, _ := user_service.GetUserId(ctx, followeeUsername)
		return followeeId
	})

	followerId := follower_id_fu.Await().(int64)
	followeeId := followee_id_fu.Await().(int64)

	if followerId == 0 || followeeId == 0 {
		fmt.Printf("Failed to find the user profile - followerUsername: %s, followeeUsername: %s\n", followerUsername, followeeUsername)
		return nil
	}
	storage := s.storage.Get()
	storage.Unfollow(ctx, followerId, followeeId)
	return nil
}
