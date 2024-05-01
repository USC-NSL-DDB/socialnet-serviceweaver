package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"SocialNetwork/client/perf"
	"SocialNetwork/shared/api"
	"SocialNetwork/shared/common"

	"github.com/ServiceWeaver/weaver/runtime/codegen"
)

const (
	NUM_THREADS            = 200
	TARGET_MOPS            = 0.1
	TOTAL_MOPS             = 1
	TIMESERIES_INTERVAL_US = 10 * 1000
)

const (
	NUM_USER         = 962
	TIMELINE_INT_MIN = 0
	TIMELINE_INT_MAX = 99
	NUM_URLS_MAX     = 2
	NUM_MEDIAS_MAX   = 2
	NUM_MENTIONS_MAX = 2
	TEXT_LEN         = 64
	URL_LEN          = 64
	CHAR_SET         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	PERCENT_USER_TIMELINE = 60
	PERCENT_HOME_TIMELINE = 30
	PERCENT_COMPOSE_POST  = 5
	PERCENT_REMOVE_POSTS  = 5
	PERCENT_FOLLOW        = 100 - PERCENT_USER_TIMELINE - PERCENT_HOME_TIMELINE - PERCENT_COMPOSE_POST - PERCENT_REMOVE_POSTS

	INTERVAL_BETWEEN_REQUESTS = 50 * time.Millisecond

	BASE_URL = "http://128.110.223.7"
)

type SingleThreadClient struct {
	rand_charset_generator rand.Rand
	rand_user_id_generator rand.Rand
	rand_int_generator     rand.Rand

	rand_max_urls_generator     rand.Rand
	rand_max_medias_generator   rand.Rand
	rand_max_mentions_generator rand.Rand

	rand_request_type_generator rand.Rand
}

func NewSingleThreadClient() *SingleThreadClient {
	client := SingleThreadClient{}
	client.Init()
	return &client
}

func (client *SingleThreadClient) Init() {
	client.rand_charset_generator = *rand.New(rand.NewSource(0))
	client.rand_user_id_generator = *rand.New(rand.NewSource(1))
	client.rand_int_generator = *rand.New(rand.NewSource(2))
	client.rand_max_urls_generator = *rand.New(rand.NewSource(3))
	client.rand_max_medias_generator = *rand.New(rand.NewSource(4))
	client.rand_max_mentions_generator = *rand.New(rand.NewSource(5))
	client.rand_request_type_generator = *rand.New(rand.NewSource(6))
}

type SocialNetworkAdapter struct{}

func (adapter *SocialNetworkAdapter) CreateGoroutineState() perf.PerfThreadState {
	return NewSingleThreadClient()
}

func (adapter *SocialNetworkAdapter) GenRequest(state perf.PerfThreadState) perf.PerfRequest {
	client := state.(*SingleThreadClient)
	rand_int := client.rand_request_type_generator.Int() % 100
	if rand_int <= PERCENT_USER_TIMELINE {
		req := &api.ReadUserTimelineRequest{}
		GenReadUserTimelineReq(req, client)
		return req
	}
	rand_int -= PERCENT_USER_TIMELINE
	if rand_int < PERCENT_HOME_TIMELINE {
		req := &api.ReadHomeTimelineRequest{}
		GenReadHomeTimelineReq(req, client)
		return req
	}

	rand_int -= PERCENT_HOME_TIMELINE
	if rand_int < PERCENT_COMPOSE_POST {
		req := &api.ComposePostRequest{}
		GenComposePostReq(req, client)
		return req
	}

	rand_int -= PERCENT_COMPOSE_POST
	if rand_int < PERCENT_REMOVE_POSTS {
		req := &api.RemovePostsRequest{}
		GenRemovePostsReq(req, client)
		return req
	}

	req := &api.FollowRequest{}
	GenFollowReq(req, client)
	return req
}

func (adapter *SocialNetworkAdapter) ServeRequest(state perf.PerfThreadState, request perf.PerfRequest) bool {
	client := state.(*SingleThreadClient)
	switch req := request.(type) {
	case *api.ReadUserTimelineRequest:
		client.SendRequest(req, BASE_URL+common.READ_USER_TIMELINE_ENDPOINT)
	case *api.ReadHomeTimelineRequest:
		client.SendRequest(req, BASE_URL+common.READ_HOME_TIMELINE_ENDPOINT)
	case *api.ComposePostRequest:
		client.SendRequest(req, BASE_URL+common.COMPOSE_POST_ENDPOINT)
	case *api.RemovePostsRequest:
		client.SendRequest(req, BASE_URL+common.REMOVE_POSTS_ENDPOINT)
	case *api.FollowRequest:
		client.SendRequest(req, BASE_URL+common.FOLLOW_ENDPOINT)
	}
	return true
}

// func (client *SingleThreadClient) GenRequest() (api.ClientRequest, string) {
// 	rand_int := client.rand_request_type_generator.Int() % 100
// 	address := BASE_URL
// 	if rand_int <= PERCENT_USER_TIMELINE {
// 		req := &api.ReadUserTimelineRequest{}
// 		GenReadUserTimelineReq(req, client)
// 		address += common.READ_USER_TIMELINE_ENDPOINT
// 		return req, address
// 	}
// 	rand_int -= PERCENT_USER_TIMELINE
// 	if rand_int < PERCENT_HOME_TIMELINE {
// 		req := &api.ReadHomeTimelineRequest{}
// 		GenReadHomeTimelineReq(req, client)
// 		address += common.READ_HOME_TIMELINE_ENDPOINT
// 		return req, address
// 	}

// 	rand_int -= PERCENT_HOME_TIMELINE
// 	if rand_int < PERCENT_COMPOSE_POST {
// 		req := &api.ComposePostRequest{}
// 		GenComposePostReq(req, client)
// 		address += common.COMPOSE_POST_ENDPOINT
// 		return req, address
// 	}

// 	rand_int -= PERCENT_COMPOSE_POST
// 	if rand_int < PERCENT_REMOVE_POSTS {
// 		req := &api.RemovePostsRequest{}
// 		GenRemovePostsReq(req, client)
// 		address += common.REMOVE_POSTS_ENDPOINT
// 		return req, address
// 	}

// 	req := &api.FollowRequest{}
// 	GenFollowReq(req, client)
// 	address += common.FOLLOW_ENDPOINT
// 	return req, address
// }

func (client *SingleThreadClient) SendRequest(req api.ClientRequest, address string) {
	data := req.Encode(codegen.NewEncoder())
	fmt.Println("Sending request to ", address, " with data: ", data)
	response, err := api.SendRequest(address, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()
}

func (client *SingleThreadClient) _gen_user_id() int64 {
	return client.rand_user_id_generator.Int63() % NUM_USER
}

func (client *SingleThreadClient) _gen_timeline_int() int {
	return client.rand_int_generator.Int()%(TIMELINE_INT_MAX-TIMELINE_INT_MIN) + TIMELINE_INT_MIN
}

func (client *SingleThreadClient) _gen_text(text_len int) string {
	text := ""
	for i := 0; i < text_len; i++ {
		text += string(CHAR_SET[client.rand_charset_generator.Int()%len(CHAR_SET)])
	}
	return text
}

func GenReadHomeTimelineReq(req *api.ReadHomeTimelineRequest, client *SingleThreadClient) {
	req.UserId = client._gen_user_id()
	req.Start = client._gen_timeline_int()
	req.Stop = req.Start + 1
}

func GenReadUserTimelineReq(req *api.ReadUserTimelineRequest, client *SingleThreadClient) {
	req.UserId = client._gen_user_id()
	req.Start = client._gen_timeline_int()
	req.Stop = req.Start + 1
}

func GenComposePostReq(req *api.ComposePostRequest, client *SingleThreadClient) {
	req.UserId = client._gen_user_id()
	req.Username = fmt.Sprintf("username_%d", req.UserId)
	req.Text = client._gen_text(TEXT_LEN)

	num_mentions := client.rand_int_generator.Int() % NUM_MENTIONS_MAX
	for i := 0; i < num_mentions; i++ {
		mention_id := client._gen_user_id()
		req.Text += fmt.Sprintf(" @username_%d", mention_id)
	}

	num_urls := client.rand_max_urls_generator.Int() % NUM_URLS_MAX
	for i := 0; i < num_urls; i++ {
		req.Text += "http://" + client._gen_text(URL_LEN)
	}

	num_medias := client.rand_max_medias_generator.Int() % NUM_MEDIAS_MAX
	req.MediaIds = make([]int64, num_medias)
	for i := 0; i < num_medias; i++ {
		req.MediaIds[i] = client.rand_user_id_generator.Int63()
		req.MediaTypes = append(req.MediaTypes, "png")
	}
	req.PostType = common.POST
}

func GenRemovePostsReq(req *api.RemovePostsRequest, client *SingleThreadClient) {
	req.UserId = client._gen_user_id()
	req.Start = 0
	req.Stop = 1
}

func GenFollowReq(req *api.FollowRequest, client *SingleThreadClient) {
	req.UserId = client._gen_user_id()
	req.FolloweeId = client._gen_user_id()
}

func GenUnfollowReq(req *api.UnfollowRequest, client *SingleThreadClient) {
	req.UserId = client._gen_user_id()
	req.FolloweeId = client._gen_user_id()
}

func GenFollowerReq(req *api.GetFollowersRequest, client *SingleThreadClient) {
	req.UserId = client._gen_user_id()
}

func writeTimeseries(perfer perf.Perf, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	timeseries := perfer.GetTimeseriesNthLats(TIMESERIES_INTERVAL_US, 99)
	for _, trace := range timeseries {
		if _, err := fmt.Fprintf(writer, "%d %d\n", trace.StartUs, trace.Duration); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to file: %v\n", err)
			return
		}
	}
	if err := writer.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush writer: %v\n", err)
	}
}

func main() {
	// client := SingleThreadClient{}
	// client.Init()
	sn_adapter := SocialNetworkAdapter{}
	perf_runner := perf.NewPerf(&sn_adapter)
	duration_us := uint64(TOTAL_MOPS / TARGET_MOPS * 1e6)
	warmup_us := duration_us
	perf_runner.Run(NUM_THREADS, TARGET_MOPS, duration_us, warmup_us, 50*1000)

	fmt.Println("real_mops, avg_lat, 50th_lat, 90th_lat, 95th_lat, 99th_lat, 99.9th_lat")
	fmt.Printf("%f %d %d %d %d %d %d\n",
		perf_runner.GetRealMops(),
		perf_runner.GetAvgLat(),
		perf_runner.GetNthLats(50),
		perf_runner.GetNthLats(90),
		perf_runner.GetNthLats(95),
		perf_runner.GetNthLats(99),
		perf_runner.GetNthLats(99.9),
	)

	writeTimeseries(*perf_runner, "timeseries.txt")

	// for {
	// 	req, address := client.GenRequest()
	// 	client.SendRequest(req, address)
	// 	time.Sleep(INTERVAL_BETWEEN_REQUESTS)
	// }
}
