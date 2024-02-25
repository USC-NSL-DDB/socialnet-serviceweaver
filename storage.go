package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
  "fmt"
)

type IStorage interface {
  // RegisterMap(context.Context, string) error
  // GetMap(context.Context, string, key) (HashMap, error)
  PutUserProfile(context.Context, string, UserProfile)
  GetUserProfile(context.Context, string) (UserProfile, bool)
  PutPost(context.Context, int64, Post)
  GetPost(context.Context, int64) (Post, bool)
  RemovePost(context.Context, int64) bool
  PutMediaData(context.Context, string, string)
  GetMediaData(context.Context, string) (string, bool)
  Follow(context.Context, int64, int64)
  Unfollow(context.Context, int64, int64)
  GetFollowers(context.Context, int64) (map[int64]bool, bool)
  GetFollowees(context.Context, int64) (map[int64]bool, bool)
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

  // mu sync.Mutex
  // data map[string]string

  usernameToUserProfileMap *HashMap[string, UserProfile]
  postIdToPostMap *HashMap[int64, Post]
  filenameToMediaDataMap *HashMap[string, string]
  useridToFollowersMap *HashMap[int64, map[int64]bool]
  useridToFolloweesMap *HashMap[int64, map[int64]bool]
}

func (s *Storage) Init(context.Context) error {
  s.usernameToUserProfileMap = NewHashMap[string, UserProfile]() 
  s.postIdToPostMap = NewHashMap[int64, Post]()
  s.filenameToMediaDataMap = NewHashMap[string, string]()
  return nil
}

func (s *Storage) PutUserProfile(_ context.Context, key string, val UserProfile) {
  s.usernameToUserProfileMap.Put(key, val)
}

func (s *Storage) GetUserProfile(_ context.Context, key string) (UserProfile, bool) {
  return s.usernameToUserProfileMap.Get(key)
}

func (s *Storage) PutPost(_ context.Context, key int64, val Post) {
  s.postIdToPostMap.Put(key, val)
}

func (s *Storage) GetPost(_ context.Context, key int64) (Post, bool) {
  return s.postIdToPostMap.Get(key)
}

func (s *Storage) RemovePost(_ context.Context, key int64) bool {
  _, exist := s.postIdToPostMap.Get(key)
  if !exist {
    return false
  }
  s.postIdToPostMap.Delete(key)
  return true
}

func (s *Storage) PutMediaData(_ context.Context, key string, val string) {
  s.filenameToMediaDataMap.Put(key, val)
}

func (s *Storage) GetMediaData(_ context.Context, key string) (string, bool) {
  return s.filenameToMediaDataMap.Get(key)
}

func (s *Storage) Follow(_ context.Context, userId int64, followeeId int64) {
  // userId follows followeeId
  followees, flag := s.useridToFolloweesMap.Get(userId)
  if !flag {
    followees = map[int64]bool{followeeId: true}
  } else {
    followees[followeeId] = true
  }
  s.useridToFolloweesMap.Put(userId, followees)

  followers, flag := s.useridToFollowersMap.Get(followeeId)
  if !flag {
    followers = map[int64]bool{userId: true}
  } else {
    followers[userId] = true
  }
  s.useridToFollowersMap.Put(followeeId, followers)
}

func (s *Storage) Unfollow(_ context.Context, userId int64, followeeId int64) {
  // userId unfollows followeeId
  followees, flag1 := s.useridToFolloweesMap.Get(userId)
  followers, flag2 := s.useridToFollowersMap.Get(followeeId)
  if !flag1 || !flag2 {
    fmt.Printf("Unfollow: userId %d or followeeId %d does not exist\n", userId, followeeId)
    return
  }

  if _, ok := followees[followeeId]; !ok {
    fmt.Printf("Unfollow: userId %d does not follow followeeId %d\n", userId, followeeId)
    return
  }

  if _, ok := followers[userId]; !ok {
    fmt.Printf("Unfollow: followeeId %d does not have userId %d as follower\n", followeeId, userId)
    return
  }

  delete(followees, followeeId)
  delete(followers, userId)
}

func (s *Storage) GetFollowers(_ context.Context, userId int64) (map[int64]bool, bool) {
  return s.useridToFollowersMap.Get(userId)
}

func (s *Storage) GetFollowees(_ context.Context, userId int64) (map[int64]bool, bool) {
  return s.useridToFolloweesMap.Get(userId)
}