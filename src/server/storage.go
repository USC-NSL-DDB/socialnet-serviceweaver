package main

import (
	"context"
	"fmt"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/btree"
)

type IStorage interface {
	// RegisterMap(context.Context, string) error
	// GetMap(context.Context, string, key) (HashMap, error)
	PutUserProfile(context.Context, string, UserProfile) error
	GetUserProfile(context.Context, string) (UserProfile, bool, error)

	PutPost(context.Context, int64, Post) error
	GetPost(context.Context, int64) (Post, bool, error)
	RemovePost(context.Context, int64) (bool, error)

	PutMediaData(context.Context, string, string) error
	GetMediaData(context.Context, string) (string, bool, error)

	PutShortenUrl(context.Context, string, string) error
	GetShortenUrl(context.Context, string) (string, bool, error)
	RemoveShortenUrl(context.Context, string) error

	Follow(context.Context, int64, int64) error
	Unfollow(context.Context, int64, int64) error
	GetFollowers(context.Context, int64) (map[int64]bool, bool, error)
	GetFollowees(context.Context, int64) (map[int64]bool, bool, error)

	PutPostTimeline(context.Context, int64, int64, int64) error
	GetPostTimeline(context.Context, int64, int, int) ([]int64, error)
	RemovePostTimeline(context.Context, int64, int64, int64) error
}

// Manually routing all request to the same replica.
// But it's not guaranteed by the documentation.
type StorageRouter struct{}

const ROUTE_KEY string = "storage"

func (StorageRouter) PutUserProfile(context.Context, string, UserProfile) string  { return ROUTE_KEY }
func (StorageRouter) GetUserProfile(context.Context, string) string               { return ROUTE_KEY }
func (StorageRouter) PutPost(context.Context, int64, Post) string                 { return ROUTE_KEY }
func (StorageRouter) GetPost(context.Context, int64) string                       { return ROUTE_KEY }
func (StorageRouter) RemovePost(context.Context, int64) string                    { return ROUTE_KEY }
func (StorageRouter) PutMediaData(context.Context, string, string) string         { return ROUTE_KEY }
func (StorageRouter) GetMediaData(context.Context, string) string                 { return ROUTE_KEY }
func (StorageRouter) PutShortenUrl(context.Context, string, string) string        { return ROUTE_KEY }
func (StorageRouter) GetShortenUrl(context.Context, string) string                { return ROUTE_KEY }
func (StorageRouter) RemoveShortenUrl(context.Context, string) string             { return ROUTE_KEY }
func (StorageRouter) Follow(context.Context, int64, int64) string                 { return ROUTE_KEY }
func (StorageRouter) Unfollow(context.Context, int64, int64) string               { return ROUTE_KEY }
func (StorageRouter) GetFollowers(context.Context, int64) string                  { return ROUTE_KEY }
func (StorageRouter) GetFollowees(context.Context, int64) string                  { return ROUTE_KEY }
func (StorageRouter) PutPostTimeline(context.Context, int64, int64, int64) string { return ROUTE_KEY }
func (StorageRouter) GetPostTimeline(context.Context, int64, int, int) string     { return ROUTE_KEY }
func (StorageRouter) RemovePostTimeline(context.Context, int64, int64, int64) string {
	return ROUTE_KEY
}

//  PutUserProfile(_ context.Context, key string) string
//  GetUserProfile(_ context.Context, key, value string) string
//  PutPost(_ context.Context, key string) string
//  GetPost(_ context.Context, key, value string) string
//  RemovePost(_ context.Context, key string) string
//  PutMediaData(_ context.Context, key, value string) string
//  GetMediaData(_ context.Context, key string) string
//  PutShortenUrl(_ context.Context, key, value string) string
//  GetShortenUrl(_ context.Context, key string) string
//  RemoveShortenUrl(_ context.Context, key, value string) string
//  Follow(_ context.Context, key string) string
//  Unfollow(_ context.Context, key, value string) string
//  GetFollowers(_ context.Context, key string) string
//  GetFollowees(_ context.Context, key, value string) string
//  PutPostTimeline(_ context.Context, key string) string
//  GetPostTimeline(_ context.Context, key, value string) string
//  RemovePostTimeline(_ context.Context, key, value string) string

// the hashtable needs to be manually managed so that it has only one instance and no replicas.
// Later on, for distributed sharded hastable, we can have more replicas as shards.
type Storage struct {
	weaver.Implements[IStorage]
	weaver.WithRouter[StorageRouter]

	// mu sync.Mutex
	// data map[string]string
	filenameToMediaDataMap   *HashMap[string, string]
	usernameToUserProfileMap *HashMap[string, UserProfile]
	postIdToPostMap          *HashMap[int64, Post]
	shortToExtendedMap       *HashMap[string, string]
	useridToFollowersMap     *HashMap[int64, *HashMap[int64, bool]]
	useridToFolloweesMap     *HashMap[int64, *HashMap[int64, bool]]

	useridToTimelineMap *HashMap[int64, *btree.BTree]
}

func (s *Storage) Init(context.Context) error {
	s.filenameToMediaDataMap = NewHashMap[string, string]()
	s.usernameToUserProfileMap = NewHashMap[string, UserProfile]()
	s.postIdToPostMap = NewHashMap[int64, Post]()
	s.shortToExtendedMap = NewHashMap[string, string]()
	s.useridToFollowersMap = NewHashMap[int64, *HashMap[int64, bool]]()
	s.useridToFolloweesMap = NewHashMap[int64, *HashMap[int64, bool]]()

	s.useridToTimelineMap = NewHashMap[int64, *btree.BTree]()
	return nil
}

func (s *Storage) PutUserProfile(_ context.Context, key string, val UserProfile) error {
	s.usernameToUserProfileMap.Put(key, val)
	return nil
}

func (s *Storage) GetUserProfile(_ context.Context, key string) (UserProfile, bool, error) {
	v, e := s.usernameToUserProfileMap.Get(key)
	return v, e, nil
}

func (s *Storage) PutPost(_ context.Context, key int64, val Post) error {
	s.postIdToPostMap.Put(key, val)
	return nil
}

func (s *Storage) GetPost(_ context.Context, key int64) (Post, bool, error) {
	v, e := s.postIdToPostMap.Get(key)
	return v, e, nil
}

func (s *Storage) RemovePost(_ context.Context, key int64) (bool, error) {
	_, exist := s.postIdToPostMap.Get(key)
	if !exist {
		return false, nil
	}
	s.postIdToPostMap.Delete(key)
	return true, nil
}

func (s *Storage) PutMediaData(_ context.Context, key string, val string) error {
	s.filenameToMediaDataMap.Put(key, val)
	return nil
}

func (s *Storage) GetMediaData(_ context.Context, key string) (string, bool, error) {
	v, e := s.filenameToMediaDataMap.Get(key)
	return v, e, nil
}

func (s *Storage) Follow(_ context.Context, userId int64, followeeId int64) error {
	// userId follows followeeId
	followees, flag := s.useridToFolloweesMap.Get(userId)
	if !flag {
		newMap := NewHashMap[int64, bool]()
		newMap.Put(followeeId, true)
		followees = newMap
	} else {
		followees.Put(followeeId, true)
	}
	s.useridToFolloweesMap.Put(userId, followees)

	followers, flag := s.useridToFollowersMap.Get(followeeId)
	if !flag {
		newMap := NewHashMap[int64, bool]()
		newMap.Put(userId, true)
		followers = newMap
	} else {
		followees.Put(userId, true)
	}
	s.useridToFollowersMap.Put(followeeId, followers)
	return nil
}

func (s *Storage) Unfollow(_ context.Context, userId int64, followeeId int64) error {
	// userId unfollows followeeId
	followees, flag1 := s.useridToFolloweesMap.Get(userId)
	followers, flag2 := s.useridToFollowersMap.Get(followeeId)
	if !flag1 || !flag2 {
		fmt.Printf("Unfollow: userId %d or followeeId %d does not exist\n", userId, followeeId)
		return nil
	}

	if _, ok := followees.Get(followeeId); !ok {
		fmt.Printf("Unfollow: userId %d does not follow followeeId %d\n", userId, followeeId)
		return nil
	}

	if _, ok := followers.Get(userId); !ok {
		fmt.Printf("Unfollow: followeeId %d does not have userId %d as follower\n", followeeId, userId)
		return nil
	}

	followees.Delete(followeeId)
	followers.Delete(userId)
	// delete(followees, followeeId)
	// delete(followers, userId)
	return nil
}

func (s *Storage) GetFollowers(_ context.Context, userId int64) (map[int64]bool, bool, error) {
	v, e := s.useridToFollowersMap.Get(userId)
	if !e {
		return nil, false, nil
	}
	return v.Clone(), e, nil
}

func (s *Storage) GetFollowees(_ context.Context, userId int64) (map[int64]bool, bool, error) {
	v, e := s.useridToFolloweesMap.Get(userId)
	if !e {
		return nil, false, nil
	}
	return v.Clone(), e, nil
}

func (s *Storage) PutShortenUrl(_ context.Context, key string, val string) error {
	s.shortToExtendedMap.Put(key, val)
	return nil
}

func (s *Storage) GetShortenUrl(_ context.Context, key string) (string, bool, error) {
	v, e := s.shortToExtendedMap.Get(key)
	return v, e, nil
}

func (s *Storage) RemoveShortenUrl(_ context.Context, key string) error {
	s.shortToExtendedMap.Delete(key)
	return nil
}

type PostTimestampPair struct {
	timestamp int64
	postId    int64
}

// Less implements btree.Item.
func (p PostTimestampPair) Less(than btree.Item) bool {
	other, ok := than.(PostTimestampPair)
	if !ok {
		return false
	}
	return p.timestamp < other.timestamp
}

func (s *Storage) PutPostTimeline(_ context.Context, userId int64, postId int64, timestamp int64) error {
	s.useridToTimelineMap.ApplyWithDefault(
		userId,
		func(k int64, v *btree.BTree, args ...interface{}) {
			timestamp := args[0].(int64)
			postId := args[1].(int64)
			v.ReplaceOrInsert(PostTimestampPair{timestamp, postId})
		},
		func(k int64) *btree.BTree {
			return btree.New(2)
		},
		timestamp, postId,
	)
	return nil
}

func (s *Storage) GetPostTimeline(_ context.Context, userId int64, start int, stop int) ([]int64, error) {
	return ApplyWithReturn(
		s.useridToTimelineMap,
		userId,
		func(k int64, v *btree.BTree, args ...interface{}) []int64 {
			start := args[0].(int)
			stop := args[1].(int)
			result := make([]int64, 0)
			v.Ascend(func(item btree.Item) bool {
				if start <= 0 {
					result = append(result, item.(PostTimestampPair).postId)
				}
				start--
				stop--
				return stop > 0
			})
			return result
		},
		start, stop,
	)
}

func (s *Storage) RemovePostTimeline(_ context.Context, userId int64, postId int64, timestamp int64) error {
	s.useridToTimelineMap.Apply(
		userId,
		func(k int64, v *btree.BTree, args ...interface{}) {
			timestamp := args[0].(int64)
			postId := args[1].(int64)
			v.Delete(PostTimestampPair{timestamp, postId})
		},
		timestamp, postId,
	)
	return nil
}
