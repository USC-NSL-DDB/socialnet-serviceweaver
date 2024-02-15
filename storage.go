package main

import (
	"context"
	"sync"

	"github.com/ServiceWeaver/weaver"
)

type IStorage interface {
  // RegisterMap(context.Context, string) error
  // GetMap(context.Context, string, key) (HashMap, error)
  PutUserProfile(context.Context, string, UserProfile)
  GetUserProfile(context.Context, string) (UserProfile, bool)
}


// Manually routing all request to the same replica.
// But it's not guaranteed by the documentation.
type StorageRouter struct {}
func (StorageRouter) Get(_ context.Context, key string) string { return "storage" }
func (StorageRouter) Put(_ context.Context, key, value string) string { return "storage" } 

// the hashtable needs to be manually managed so that it has only one instance and no replicas.
// Later on, for distributed sharded hastable, we can have more replicas as shards.
type Storage struct {
  weaver.Implements[IStorage]
  weaver.WithRouter[StorageRouter]

  mu sync.Mutex
  data map[string]string

  usernameToUserProfileMap *HashMap[string, UserProfile]
}

func (s *Storage) Init(context.Context) error {
  s.usernameToUserProfileMap = NewHashMap[string, UserProfile]() 
  return nil
}

func (s *Storage) PutUserProfile(_ context.Context, key string, val UserProfile) {
  s.usernameToUserProfileMap.Put(key, val)
}

func (s *Storage) GetUserProfile(_ context.Context, key string) (UserProfile, bool) {
  return s.usernameToUserProfileMap.Get(key)
}

