package main

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
)

type IHomeTimelineService interface {
	ReadHomeTimeline(context.Context, int64, int, int) []Post
	WriteHomeTimeline(context.Context, int64, int64, int64, []int64)
	RemovePost(context.Context, int64, int64, int64)
}

type HomeTimelineService struct {
	weaver.Implements[IHomeTimelineService]

	postStorageService weaver.Ref[PostStorageService]
	socialGraphService weaver.Ref[SocialGraphService]
	storage            weaver.Ref[Storage]
}

func (hts *HomeTimelineService) ReadHomeTimeline(ctx context.Context, userId int64, start int, stop int) []Post {
	if stop <= start || start < 0 {
		return make([]Post, 0)
	}
	storage := hts.storage.Get()
	postStorageService := hts.postStorageService.Get()

	postIds := storage.GetPostTimeline(ctx, userId, start, stop)
	return postStorageService.ReadPosts(ctx, postIds)
}

func (hts *HomeTimelineService) WriteHomeTimeline(ctx context.Context, postId int64, userId int64, timestamp int64, userMentionIds []int64) {
	storage := hts.storage.Get()
	socialGraphService := hts.socialGraphService.Get()
	ids := socialGraphService.GetFollowers(ctx, userId)
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(ctx context.Context, id, postId, timestamp int64) {
			defer wg.Done()
			storage.PutPostTimeline(ctx, id, postId, timestamp)
		}(ctx, id, postId, timestamp)
	}
	wg.Wait()
}

func (hts *HomeTimelineService) RemovePost(ctx context.Context, userId int64, postId int64, timestamp int64) {
	storage := hts.storage.Get()
	storage.RemovePostTimeline(ctx, userId, postId, timestamp)
}
