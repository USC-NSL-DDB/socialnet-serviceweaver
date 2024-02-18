package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"fmt"
)

type PostStorageServicer interface {
	RemovePost(context.Context, int64) bool
	ReadPost(context.Context, int64) Post
	StorePost(context.Context, Post)
	ReadPosts(context.Context, []int64) []Post
}


type PostStorageService struct {
	weaver.Implements[PostStorageServicer]
	storage weaver.Ref[Storage]
	
}


func (pss *PostStorageService) StorePost(ctx context.Context, post Post) {
	storage := pss.storage.Get()
	storage.PutPost(ctx, post.post_id, post)
}

func (pss *PostStorageService) ReadPost(ctx context.Context, postId int64) Post {
	storage := pss.storage.Get()
	post, exist := storage.GetPost(ctx, postId)
	if !exist {
		fmt.Printf("Failed to find the post - post_id: %d\n", postId)
		return Post{}
	}
	return post
}


func (pss *PostStorageService) ReadPosts(ctx context.Context, postIds []int64) []Post {
	storage := pss.storage.Get()
	posts := make([]Post, 0)
	for _, postId := range postIds {
		post, exist := storage.GetPost(ctx, postId)
		if !exist {
			fmt.Printf("Failed to find the post - post_id: %d\n", postId)
			post = Post{}
		}
		posts = append(posts, post)
	}
	return posts
}

func (pss *PostStorageService) RemovePost(ctx context.Context, postId int64) bool {
	storage := pss.storage.Get()
	return storage.RemovePost(ctx, postId)
}
