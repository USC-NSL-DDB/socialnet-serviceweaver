package main

import (
	"context"

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
	CompostPost(context.Context, string, int64, string, []int64, []string)
	Login(context.Context, string, string) string
	RegisterUser(context.Context, string, string, string, string)
	RegisterUserWithId(context.Context, string, string, string, string, int64)
	ReadUserTimeline(context.Context, int64, int, int)
	// Reverse(context.Context, string) (string, error)
}

type BackendService struct {
	weaver.Implements[BackendServicer]

	userService         weaver.Ref[UserServicer]
	userTimelineService weaver.Ref[UserTimelineService]
	socialGraphService  weaver.Ref[SocialGraphService]
	postStorageService  weaver.Ref[PostStorageService]
	homeTimelineService weaver.Ref[HomeTimelineService]
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

func (bs *BackendService) RemovePosts(user_id int64, start, top int) {
	// run UserTimelineService
	// run SocialGraphService
	// run PostStorageService
	// run HomeTimelineService
	// run UrlShortenService
}

func (bs *BackendService) CompostPost(
	ctx context.Context,
	username string,
	user_id int64,
	text string,
	media_ids []int64,
	media_types []string,
) {
	// run TextService
	// run UniqueIdService
	// run MediaService
	// run UserSerivce
	// run UserTimelineService
	// run HomeTimelineService
}

func (bs *BackendService) ReadUserTimeline(
	ctx context.Context,
	user_id int64,
	start, stop int,
) {
	// run ReadUserTimelineService
}

// func (r *reverser) Reverse(_ context.Context, s string) (string, error) {
//     runes := []rune(s)
//     n := len(runes)
//     for i := 0; i < n/2; i++ {
//         runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
//     }
//     return string(runes), nil
// }
