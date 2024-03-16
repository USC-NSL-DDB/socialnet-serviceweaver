package main

import (
	"context"

	"fmt"

	"github.com/ServiceWeaver/weaver"
)

type MediaStorageServicer interface {
	UploadMedia(context.Context, string, string)
	GetMedia(context.Context, string) (string, error)
}

type MediaStorageService struct {
	weaver.Implements[MediaStorageServicer]
	storage weaver.Ref[Storage]
}

func (m *MediaStorageService) UploadMedia(ctx context.Context, filename string, data string) {
	storage := m.storage.Get()
	storage.PutMediaData(ctx, filename, data)
}

func (m *MediaStorageService) GetMedia(ctx context.Context, filename string) (string, error) {
	storage := m.storage.Get()
	data, ok, _ := storage.GetMediaData(ctx, filename)
	if !ok {
		fmt.Printf("Failed to find the media - filename: %s\n", filename)
		return "", nil
	}
	return data, nil
}
