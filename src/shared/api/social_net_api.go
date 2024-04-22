package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"SocialNetwork/shared/common"

	"github.com/ServiceWeaver/weaver/runtime/codegen"
)

var clientPool = sync.Pool{
	New: func() interface{} {
		return &http.Client{}
	},
}

func getClient() *http.Client {
	return clientPool.Get().(*http.Client)
}

func releaseClient(client *http.Client) {
	clientPool.Put(client)
}

const (
	BASE_PORT = "49555"
	BASE_URL  = "http://localhost:" + BASE_PORT
)

// var client = &http.Client{}

func SendRequest(address string, data []byte) (*http.Response, error) {
	// fmt.Println("Sending to addr:", address)
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("[SendRequest] Error:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/custom")
	client := getClient()
	resp, err := client.Do(req)
	releaseClient(client)
	return resp, err
}

func InitRequest(endpoint string, data []byte) (*http.Response, error) {
	address := BASE_URL + endpoint
	return SendRequest(address, data)
}

func EncodeData(action func(*codegen.Encoder)) []byte {
	enc := codegen.NewEncoder()
	action(enc)
	return enc.Data()
}

func DecodeData(response *http.Response, action func(*codegen.Decoder)) {
	defer response.Body.Close()
	// fmt.Println("Status Code:", response.StatusCode)

	resp_body, _ := io.ReadAll(response.Body)
	dec := codegen.NewDecoder(resp_body)
	action(dec)
}

func send_request_wrapper(full_addr string, req EncodableRequest) (*http.Response, error) {
	return SendRequest(full_addr, req.Encode(codegen.NewEncoder()))
}

func RegisterUser(addr string, req *RegisterUserRequest) {
	resp, err := send_request_wrapper(addr+common.REGISTER_USER_ENDPOINT, req)
	if err != nil {
		fmt.Println("[RegisterUser] Error:", err)
		return
	}
	defer resp.Body.Close()
}

func RegisterUserWithId(addr string, req *RegisterUserWithIdRequest) {
	resp, err := send_request_wrapper(addr+common.REGISTER_USER_WITH_ID_ENDPOINT, req)
	if err != nil {
		fmt.Println("[RegisterUserWithId] Error:", err)
		return
	}
	defer resp.Body.Close()
}

func ReadHomeTimeline(addr string, req *ReadHomeTimelineRequest) {
	resp, err := send_request_wrapper(addr+common.READ_HOME_TIMELINE_ENDPOINT, req)
	if err != nil {
		fmt.Println("[ReadHomeTimeline] Error:", err)
		return
	}
	defer resp.Body.Close()
}

func ReadUserTimeline(addr string, req *ReadUserTimelineRequest) {
	resp, err := send_request_wrapper(addr+common.READ_USER_TIMELINE_ENDPOINT, req)
	if err != nil {
		fmt.Println("[ReadUserTimeline] Error:", err)
		return
	}
	defer resp.Body.Close()
}

func RemovePosts(addr string, req *RemovePostsRequest) {
	resp, err := send_request_wrapper(addr+common.REMOVE_POSTS_ENDPOINT, req)
	if err != nil {
		fmt.Println("[RemovePosts] Error:", err)
		return
	}
	defer resp.Body.Close()
}

func ComposePost(addr string, req *ComposePostRequest) {
	resp, err := send_request_wrapper(addr+common.COMPOSE_POST_ENDPOINT, req)
	if err != nil {
		fmt.Println("[ComposePost] Error:", err)
		return
	}
	defer resp.Body.Close()
}

func Login(addr string, req *LoginRequest) {
	resp, err := send_request_wrapper(addr+common.LOGIN_ENDPOINT, req)
	if err != nil {
		fmt.Println("[Login] Error:", err)
		return
	}
	defer resp.Body.Close()
}

func Follow(addr string, req *FollowRequest) {
	resp, err := send_request_wrapper(addr+common.FOLLOW_ENDPOINT, req)
	if err != nil {
		fmt.Println("[Follow] Error:", err)
		return
	}
	defer resp.Body.Close()
}

func FollowWithUsername(addr string, req *FollowWithUsernameRequest) {
	resp, err := send_request_wrapper(addr+common.FOLLOW_WITH_USERNAME_ENDPOINT, req)
	if err != nil {
		fmt.Println("[FollowWithUsername] Error:", err)
		return
	}
	defer resp.Body.Close()
}

func Unfollow(addr string, req *UnfollowRequest) {
	resp, err := send_request_wrapper(addr+common.UNFOLLOW_ENDPOINT, req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}

func UnfollowWithUsername(addr string, req *UnfollowWithUsernameRequest) {
	resp, err := send_request_wrapper(addr+common.UNFOLLOW_WITH_USERNAME_ENDPOINT, req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}

func GetFollowers(addr string, req *GetFollowersRequest) {
	resp, err := send_request_wrapper(addr+common.GET_FOLLOWERS_ENDPOINT, req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}

func GetFollowees(addr string, req *GetFolloweesRequest) {
	resp, err := send_request_wrapper(common.GET_FOLLOWEES_ENDPOINT, req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}

func UploadMedia(addr string, req *UploadMediaRequest) {
	resp, err := send_request_wrapper(common.UPLOAD_MEDIA_ENDPOINT, req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}

func GetMedia(addr string, req *GetMediaRequest) {
	resp, err := send_request_wrapper(addr+common.GET_MEDIA_ENDPOINT, req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}
