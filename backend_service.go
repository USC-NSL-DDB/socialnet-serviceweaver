package main

import (
	"context"
	"time"

	"github.com/ServiceWeaver/weaver"
)

type BackendServicer interface {
	//   void RemovePosts(int64_t user_id, int start, int stop);
	//   void ComposePost(const std::string &username, int64_t user_id,
	//                    const std::string &text,
	//                    const std::vector<int64_t> &media_ids,
	//                    const std::vector<std::string> &media_types,
	//                    PostType::type post_type) override;
	//   void ReadUserTimeline(std::vector<Post> &, int64_t, int, int) override;
	//   void Login(std::string &_return, const std::string &username,
	//              const std::string &password) override;
	//   void RegisterUser(const std::string &first_name, const std::string &last_name,
	//                     const std::string &username,
	//                     const std::string &password) override;
	//   void RegisterUserWithId(const std::string &first_name,
	//                           const std::string &last_name,
	//                           const std::string &username,
	//                           const std::string &password,
	//                           const int64_t user_id) override;
	//   void GetFollowers(std::vector<int64_t> &_return,
	//                     const int64_t user_id) override;
	//   void Unfollow(const int64_t user_id, const int64_t followee_id) override;
	//   void UnfollowWithUsername(const std::string &user_usernmae,
	//                             const std::string &followee_username) override;
	//   void Follow(const int64_t user_id, const int64_t followee_id) override;
	//   void FollowWithUsername(const std::string &user_usernmae,
	//                           const std::string &followee_username) override;
	//   void GetFollowees(std::vector<int64_t> &_return,
	//                     const int64_t user_id) override;
	//   void ReadHomeTimeline(std::vector<Post> &_return, const int64_t user_id,
	//                         const int32_t start, const int32_t stop) override;
	//   void UploadMedia(const std::string &filename,
	//                    const std::string &data) override;
	//   void GetMedia(std::string &_return, const std::string &filename) override;
	RemovePosts(context.Context, int64, int)
	CompostPost(context.Context, string, int64, string, []int64, []string, PostType)
	Login(context.Context, string, string) string
	RegisterUser(context.Context, string, string, string, string)
	RegisterUserWithId(context.Context, string, string, string, string, int64)
	ReadUserTimeline(context.Context, int64, int, int) []Post
	GetFollowers(context.Context, int64) []int64
	Unfollow(context.Context, int64, int64)
	UnfollowWithUsername(context.Context, string, string)
	GetFollowees(context.Context, int64) []int64
	ReadHomeTimeline(context.Context, int64, int, int) []Post
	UploadMedia(context.Context, string, string)
	GetMedia(context.Context, string) string
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
	mediaService        weaver.Ref[MediaService]
}

func (bs *BackendService) Login(ctx context.Context, username string, password string) string {
	variant, err := bs.userService.Get().Login(ctx, username, password)
	if err != nil {
		panic(err)
	}
	return variant
}

func (bs *BackendService) RegisterUser(
	ctx context.Context,
	first_name,
	last_name,
	username,
	password string,
) {
	// run UserService
	bs.userService.Get().RegisterUser(ctx, first_name, last_name, username, password)
}

func (bs *BackendService) RegisterUserWithId(
	ctx context.Context,
	first_name,
	last_name,
	username,
	password string,
	user_id int64,
) {
	// run UserService
	bs.userService.Get().RegisterUserWithId(ctx, first_name, last_name, username, password, user_id)
}

func (bs *BackendService) RemovePosts(ctx context.Context, user_id int64, start, top int) {
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
		return utls.ReadUserTimeline(ctx, user_id, start, top)
	})

	followers_fu := AsyncExec(func() interface{} {
		return sgs.GetFollowers(ctx, user_id)
	})

	posts := posts_fu.Await().([]Post)
	followers := followers_fu.Await().([]int64)

	remove_posts_fus := make([]Future, 0)
	remove_from_timeline_fus := make([]Future, 0)
	remove_short_url_fus := make([]Future, 0)

	for _, post := range posts {
		remove_posts_fus = append(remove_posts_fus, AsyncExec(func() interface{} {
			return pss.RemovePost(ctx, post.post_id)
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
}

func (bs *BackendService) CompostPost(
	ctx context.Context,
	username string,
	user_id int64,
	text string,
	media_ids []int64,
	media_types []string,
	post_type PostType,
) {
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

	text_fu := AsyncExec(func() interface{} { return text_service.ComposeText(ctx, text) })
	unique_id_fu := AsyncExec(func() interface{} { return unique_id_service.ComposeUniqueId(ctx, post_type) })
	medias_fu := AsyncExec(func() interface{} { return media_service.ComposeMedia(ctx, media_types, media_ids) })
	creator_fu := AsyncExec(func() interface{} { return us.ComposeCreatorWithUserId(ctx, user_id, username) })

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
}

func (bs *BackendService) ReadUserTimeline(
	ctx context.Context,
	user_id int64,
	start, stop int,
) []Post {
	// run ReadUserTimelineService
	utls := bs.userTimelineService.Get()
	return utls.ReadUserTimeline(ctx, user_id, start, stop)
}

func (bs *BackendService) GetFollowers(context.Context, int64) []int64 {
  
}

func (bs *BackendService) Unfollow(context.Context, int64, int64) {

}

func (bs *BackendService) UnfollowWithUsername(context.Context, string, string) {

}

func (bs *BackendService) GetFollowees(context.Context, int64) []int64 {

}

func (bs *BackendService) ReadHomeTimeline(context.Context, int64, int, int) []Post {

}

func (bs *BackendService) UploadMedia(context.Context, string, string) {

}

func (bs *BackendService) GetMedia(context.Context, string) string {

}
