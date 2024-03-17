package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type IMediaService interface {
	ComposeMedia(context.Context, []string, []int64) ([]Media, error)
}

type MediaService struct {
	weaver.Implements[IMediaService]
	// storage weaver.Ref[Storage]
}

func (ms *MediaService) ComposeMedia(ctx context.Context, mediaTypes []string, mediaIds []int64) ([]Media, error) {
	media := make([]Media, 0)
	for i := 0; i < len(mediaIds); i++ {
		oneMedia := Media{
			mediaId:   mediaIds[i],
			mediaType: mediaTypes[i],
		}
		media = append(media, oneMedia)
	}
	return media, nil
}
