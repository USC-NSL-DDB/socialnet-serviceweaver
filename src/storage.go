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
	PutUserProfile(context.Context, string, UserProfile)
	GetUserProfile(context.Context, string) (UserProfile, bool, error)

	PutPost(context.Context, int64, Post)
	GetPost(context.Context, int64) (Post, bool)
	RemovePost(context.Context, int64) (bool, error)

	PutMediaData(context.Context, string, string)
	GetMediaData(context.Context, string) (string, bool, error)

	PutShortenUrl(context.Context, string, string)
	GetShortenUrl(context.Context, string) (string, bool, error)
	RemoveShortenUrl(context.Context, string)

	Follow(context.Context, int64, int64)
	Unfollow(context.Context, int64, int64)
	GetFollowers(context.Context, int64) (map[int64]bool, bool, error)
	GetFollowees(context.Context, int64) (map[int64]bool, bool, error)

	PutPostTimeline(context.Context, int64, int64, int64)
	GetPostTimeline(context.Context, int64, int, int) ([]int64, error)
	RemovePostTimeline(context.Context, int64, int64, int64)
}

// Manually routing all request to the same replica.
// But it's not guaranteed by the documentation.
type StorageRouter struct{}

func (StorageRouter) Get(_ context.Context, key string) string        { return "storage" }
func (StorageRouter) Put(_ context.Context, key, value string) string { return "storage" }

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
	useridToFollowersMap     *HashMap[int64, map[int64]bool]
	useridToFolloweesMap     *HashMap[int64, map[int64]bool]

	useridToTimelineMap *HashMap[int64, *btree.BTree]
}

func (s *Storage) Init(context.Context) error {
	s.filenameToMediaDataMap = NewHashMap[string, string]()
	s.usernameToUserProfileMap = NewHashMap[string, UserProfile]()
	s.postIdToPostMap = NewHashMap[int64, Post]()
	s.shortToExtendedMap = NewHashMap[string, string]()
	s.useridToFollowersMap = NewHashMap[int64, map[int64]bool]()
	s.useridToFolloweesMap = NewHashMap[int64, map[int64]bool]()

	s.useridToTimelineMap = NewHashMap[int64, *btree.BTree]()
	return nil
}

func (s *Storage) PutUserProfile(_ context.Context, key string, val UserProfile) {
	s.usernameToUserProfileMap.Put(key, val)
}

func (s *Storage) GetUserProfile(_ context.Context, key string) (UserProfile, bool, error) {
	v, e := s.usernameToUserProfileMap.Get(key)
	return v, e, nil
}

func (s *Storage) PutPost(_ context.Context, key int64, val Post) {
	s.postIdToPostMap.Put(key, val)
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

func (s *Storage) PutMediaData(_ context.Context, key string, val string) {
	s.filenameToMediaDataMap.Put(key, val)
}

func (s *Storage) GetMediaData(_ context.Context, key string) (string, bool, error) {
	v, e := s.filenameToMediaDataMap.Get(key)
	return v, e, nil
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

func (s *Storage) GetFollowers(_ context.Context, userId int64) (map[int64]bool, bool, error) {
	v, e := s.useridToFollowersMap.Get(userId)
	return v, e, nil
}

func (s *Storage) GetFollowees(_ context.Context, userId int64) (map[int64]bool, bool, error) {
	v, e := s.useridToFolloweesMap.Get(userId)
	return v, e, nil
}

func (s *Storage) PutShortenUrl(_ context.Context, key string, val string) {
	s.shortToExtendedMap.Put(key, val)
}

func (s *Storage) GetShortenUrl(_ context.Context, key string) (string, bool, error) {
	v, e := s.shortToExtendedMap.Get(key)
	return v, e, nil
}

func (s *Storage) RemoveShortenUrl(_ context.Context, key string) {
	s.shortToExtendedMap.Delete(key)
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

func (s *Storage) PutPostTimeline(_ context.Context, userId int64, postId int64, timestamp int64) {
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
	), nil
}

func (s *Storage) RemovePostTimeline(_ context.Context, userId int64, postId int64, timestamp int64) {
	s.useridToTimelineMap.Apply(
		userId,
		func(k int64, v *btree.BTree, args ...interface{}) {
			timestamp := args[0].(int64)
			postId := args[1].(int64)
			v.Delete(PostTimestampPair{timestamp, postId})
		},
		timestamp, postId,
	)
}
