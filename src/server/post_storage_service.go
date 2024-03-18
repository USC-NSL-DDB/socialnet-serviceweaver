package main

import (
	"context"

	"fmt"

	"github.com/ServiceWeaver/weaver"
)

type PostStorageServicer interface {
	RemovePost(context.Context, int64) (bool, error)
	ReadPost(context.Context, int64) (Post, error)
	StorePost(context.Context, Post) error
	ReadPosts(context.Context, []int64) ([]Post, error)
}

type PostStorageService struct {
	weaver.Implements[PostStorageServicer]
	storage weaver.Ref[IStorage]
}

func (pss *PostStorageService) StorePost(ctx context.Context, post Post) error {
	storage := pss.storage.Get()
	storage.PutPost(ctx, post.Post_id, post)
	return nil
}

func (pss *PostStorageService) ReadPost(ctx context.Context, postId int64) (Post, error) {
	storage := pss.storage.Get()
	post, exist, _ := storage.GetPost(ctx, postId)
	if !exist {
		fmt.Printf("Failed to find the post - post_id: %d\n", postId)
		return Post{}, nil
	}
	return post, nil
}

func (pss *PostStorageService) ReadPosts(ctx context.Context, postIds []int64) ([]Post, error) {
	storage := pss.storage.Get()
	posts := make([]Post, 0)
	for _, postId := range postIds {
		post, exist, _ := storage.GetPost(ctx, postId)
		if !exist {
			fmt.Printf("Failed to find the post - post_id: %d\n", postId)
			post = Post{}
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (pss *PostStorageService) RemovePost(ctx context.Context, postId int64) (bool, error) {
	storage := pss.storage.Get()
	return storage.RemovePost(ctx, postId)
}
