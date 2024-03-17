package main

import (
	"context"
	"time"

	"github.com/ServiceWeaver/weaver"
)

type BackendServicer interface {
	RemovePosts(context.Context, int64, int, int) error
	CompostPost(context.Context, string, int64, string, []int64, []string, PostType) error
	Login(context.Context, string, string) (string, error)
	RegisterUser(context.Context, string, string, string, string) error
	RegisterUserWithId(context.Context, string, string, string, string, int64) error
	ReadUserTimeline(context.Context, int64, int, int) ([]Post, error)
	GetFollowers(context.Context, int64) ([]int64, error)
	Unfollow(context.Context, int64, int64) error
	UnfollowWithUsername(context.Context, string, string) error
	Follow(context.Context, int64, int64) error
	FollowWithUsername(context.Context, string, string) error
	GetFollowees(context.Context, int64) ([]int64, error)
	ReadHomeTimeline(context.Context, int64, int, int) ([]Post, error)
	UploadMedia(context.Context, string, string) error
	GetMedia(context.Context, string) (string, error)
}

type BackendService struct {
	weaver.Implements[BackendServicer]

	userService         weaver.Ref[UserServicer]
	userTimelineService weaver.Ref[UserTimelineService]
	socialGraphService  weaver.Ref[SocialGraphService]
	postStorageService  weaver.Ref[PostStorageService]
	homeTimelineService weaver.Ref[HomeTimelineService]
	urlShortenService   weaver.Ref[UrlShortenService]
	textService         weaver.Ref[TextService]
	uniqueIdService     weaver.Ref[UniqueIdService]
	mediaStorageService weaver.Ref[MediaStorageService]
	mediaService        weaver.Ref[MediaService]
}

func (bs *BackendService) Login(ctx context.Context, username string, password string) (string, error) {
	variant, err := bs.userService.Get().Login(ctx, username, password)
	if err != nil {
		panic(err)
	}
	return variant, nil
}

func (bs *BackendService) RegisterUser(
	ctx context.Context,
	first_name,
	last_name,
	username,
	password string,
) error {
	// run UserService
	bs.userService.Get().RegisterUser(ctx, first_name, last_name, username, password)
	return nil
}

func (bs *BackendService) RegisterUserWithId(
	ctx context.Context,
	first_name,
	last_name,
	username,
	password string,
	user_id int64,
) error {
	// run UserService
	bs.userService.Get().RegisterUserWithId(ctx, first_name, last_name, username, password, user_id)
	return nil
}

func (bs *BackendService) RemovePosts(ctx context.Context, user_id int64, start, top int) error {
	// run UserTimelineService
	// run SocialGraphService
	// run PostStorageService
	// run HomeTimelineService
	// run UrlShortenService
	utls := bs.userTimelineService.Get()
	sgs := bs.socialGraphService.Get()
	pss := bs.postStorageService.Get()
	htls := bs.homeTimelineService.Get()
	uss := bs.urlShortenService.Get()

	posts_fu := AsyncExec(func() interface{} {
		r, _ := utls.ReadUserTimeline(ctx, user_id, start, top)
		return r
	})

	followers_fu := AsyncExec(func() interface{} {
		r, _ := sgs.GetFollowers(ctx, user_id)
		return r
	})

	posts := posts_fu.Await().([]Post)
	followers := followers_fu.Await().([]int64)

	remove_posts_fus := make([]Future, 0)
	remove_from_timeline_fus := make([]Future, 0)
	remove_short_url_fus := make([]Future, 0)

	for _, post := range posts {
		remove_posts_fus = append(remove_posts_fus, AsyncExec(func() interface{} {
			result, _ := pss.RemovePost(ctx, post.post_id)
			return result
		}).(Future))

		remove_from_timeline_fus = append(remove_from_timeline_fus, AsyncExec(func() interface{} {
			utls.RemovePost(ctx, user_id, post.post_id, post.timestamp)
			return nil
		}).(Future))

		for _, mention := range post.user_mentions {
			remove_short_url_fus = append(remove_short_url_fus, AsyncExec(func() interface{} {
				htls.RemovePost(ctx, mention.userId, post.post_id, post.timestamp)
				return nil
			}).(Future))
		}

		for _, user_id := range followers {
			remove_from_timeline_fus = append(remove_from_timeline_fus, AsyncExec(func() interface{} {
				utls.RemovePost(ctx, user_id, post.post_id, post.timestamp)
				return nil
			}).(Future))
		}

		shortened_urls := make([]string, 0)
		for _, url := range post.urls {
			shortened_urls = append(shortened_urls, url.shortenedUrl)
		}

		remove_short_url_fus = append(remove_short_url_fus, AsyncExec(func() interface{} {
			uss.RemoveUrls(ctx, shortened_urls)
			return nil
		}).(Future))
	}

	// This blocking call is not necessary in the original code
	for _, fu := range remove_posts_fus {
		fu.Await()
	}
	for _, fu := range remove_from_timeline_fus {
		fu.Await()
	}
	for _, fu := range remove_short_url_fus {
		fu.Await()
	}

	return nil
}

func (bs *BackendService) CompostPost(
	ctx context.Context,
	username string,
	user_id int64,
	text string,
	media_ids []int64,
	media_types []string,
	post_type PostType,
) error {
	// run TextService
	// run UniqueIdService
	// run MediaService
	// run UserSerivce
	// run UserTimelineService
	// run HomeTimelineService
	// run PostStorageService
	text_service := bs.textService.Get()
	unique_id_service := bs.uniqueIdService.Get()
	media_service := bs.mediaService.Get()
	us := bs.userService.Get()
	utls := bs.userTimelineService.Get()
	htls := bs.homeTimelineService.Get()
	post_storage_service := bs.postStorageService.Get()

	text_fu := AsyncExec(func() interface{} {
		r, _ := text_service.ComposeText(ctx, text)
		return r
	})
	unique_id_fu := AsyncExec(func() interface{} {
		r, _ := unique_id_service.ComposeUniqueId(ctx, post_type)
		return r
	})
	medias_fu := AsyncExec(func() interface{} {
		r, _ := media_service.ComposeMedia(ctx, media_types, media_ids)
		return r
	})
	creator_fu := AsyncExec(func() interface{} {
		r, _ := us.ComposeCreatorWithUserId(ctx, user_id, username)
		return r
	})

	timestamp := time.Now().Unix()
	unique_id := unique_id_fu.Await().(int64)

	write_user_timeline_fu := AsyncExec(func() interface{} {
		utls.WriteUserTimeline(ctx, unique_id, user_id, timestamp)
		return nil
	})

	text_service_return := text_fu.Await().(TextServiceReturn)
	user_mention_ids := make([]int64, 0)
	for _, item := range text_service_return.user_mentions {
		user_mention_ids = append(user_mention_ids, item.userId)
	}
	write_home_timeline_fu := AsyncExec(func() interface{} {
		htls.WriteHomeTimeline(ctx, unique_id, user_id, timestamp, user_mention_ids)
		return nil
	})

	post := Post{
		post_id:       unique_id,
		creator:       creator_fu.Await().(Creator),
		req_id:        0,
		text:          text_service_return.text,
		user_mentions: text_service_return.user_mentions,
		media:         medias_fu.Await().([]Media),
		urls:          text_service_return.urls,
		timestamp:     timestamp,
		post_type:     post_type,
	}

	post_fu := AsyncExec(func() interface{} {
		post_storage_service.StorePost(ctx, post)
		return nil
	})
	write_user_timeline_fu.Await()
	post_fu.Await()
	write_home_timeline_fu.Await()
	return nil
}

func (bs *BackendService) ReadUserTimeline(
	ctx context.Context,
	user_id int64,
	start, stop int,
) ([]Post, error) {
	// run ReadUserTimelineService
	utls := bs.userTimelineService.Get()
	return utls.ReadUserTimeline(ctx, user_id, start, stop)
}

func (bs *BackendService) GetFollowers(ctx context.Context, user_id int64) ([]int64, error) {
	sgs := bs.socialGraphService.Get()
	return sgs.GetFollowers(ctx, user_id)
}

func (bs *BackendService) Unfollow(ctx context.Context, user_id int64, followee_id int64) error {
	sgs := bs.socialGraphService.Get()
	sgs.Unfollow(ctx, user_id, followee_id)
	return nil
}

func (bs *BackendService) UnfollowWithUsername(ctx context.Context, user_username string, followee_username string) error {
	sgs := bs.socialGraphService.Get()
	sgs.UnfollowWithUsername(ctx, user_username, followee_username)
	return nil
}

func (bs *BackendService) Follow(ctx context.Context, user_id int64, followee_id int64) error {
	sgs := bs.socialGraphService.Get()
	sgs.Follow(ctx, user_id, followee_id)
	return nil
}

func (bs *BackendService) FollowWithUsername(ctx context.Context, user_username string, followee_username string) error {
	sgs := bs.socialGraphService.Get()
	sgs.FollowWithUsername(ctx, user_username, followee_username)
	return nil
}

func (bs *BackendService) GetFollowees(ctx context.Context, user_id int64) ([]int64, error) {
	sgs := bs.socialGraphService.Get()
	return sgs.GetFollowees(ctx, user_id)
}

func (bs *BackendService) ReadHomeTimeline(ctx context.Context, user_id int64, start int, stop int) ([]Post, error) {
	htls := bs.homeTimelineService.Get()
	return htls.ReadHomeTimeline(ctx, user_id, start, stop)
}

func (bs *BackendService) UploadMedia(ctx context.Context, filename string, data string) error {
	mss := bs.mediaStorageService.Get()
	mss.UploadMedia(ctx, filename, data)
	return nil
}

func (bs *BackendService) GetMedia(ctx context.Context, filename string) (string, error) {
	mss := bs.mediaStorageService.Get()
	return mss.GetMedia(ctx, filename)
}
