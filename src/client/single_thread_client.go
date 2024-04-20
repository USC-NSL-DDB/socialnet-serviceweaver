package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	. "SocialNetwork/shared/common"

	"github.com/ServiceWeaver/weaver/runtime/codegen"
)

const (
	NUM_USER = 962
	TIMELINE_INT_MIN = 0
	TIMELINE_INT_MAX = 99
	NUM_URLS_MAX = 2
	NUM_MEDIAS_MAX = 2
	NUM_MENTIONS_MAX = 2
	TEXT_LEN = 64
	URL_LEN = 64
	CHAR_SET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	PERCENT_USER_TIMELINE = 60
	PERCENT_HOME_TIMELINE = 30
	PERCENT_COMPOSE_POST = 5
	PERCENT_REMOVE_POSTS = 5
	PERCENT_FOLLOW = 100 - PERCENT_USER_TIMELINE - PERCENT_HOME_TIMELINE - PERCENT_COMPOSE_POST - PERCENT_REMOVE_POSTS

	INTERVAL_BETWEEN_REQUESTS = 50 * time.Millisecond
)


type ClientRequest interface {
	encode(*codegen.Encoder) []byte
	generate(client *SingleThreadClient)
}


type ClientResponse interface {
	decode(*http.Response) []byte
	display()
}


type SingleThreadClient struct {
	rand_charset_generator rand.Rand
	rand_user_id_generator rand.Rand
	rand_int_generator rand.Rand

	rand_max_urls_generator rand.Rand
	rand_max_medias_generator rand.Rand
	rand_max_mentions_generator rand.Rand

	rand_request_type_generator rand.Rand
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

func (client *SingleThreadClient) GenRequest() (ClientRequest, string) {
	rand_int := client.rand_request_type_generator.Int() % 100
	address := "http://localhost:12345/"
	if rand_int <= PERCENT_USER_TIMELINE {
		req := &UserTimelineRequest{}
		req.generate(client)
		address += READ_USER_TIMELINE_ENDPOINT
		return req, address
	}
	rand_int -= PERCENT_USER_TIMELINE
	if rand_int < PERCENT_HOME_TIMELINE {
		req := &HomeTimelineRequest{}
		req.generate(client)
		address += READ_HOME_TIMELINE_ENDPOINT
		return req, address
	}

	rand_int -= PERCENT_HOME_TIMELINE
	if rand_int < PERCENT_COMPOSE_POST {
		req := &ComposePostRequest{}
		req.generate(client)
		address += COMPOSE_POST_ENDPOINT
		return req, address
	}

	rand_int -= PERCENT_COMPOSE_POST
	if rand_int < PERCENT_REMOVE_POSTS {
		req := &RemovePostsRequest{}
		req.generate(client)
		address += REMOVE_POSTS_ENDPOINT
		return req, address
	}

	req := &FollowRequest{}
	req.generate(client)
	address += FOLLOW_ENDPOINT
	return req, address
}

func (client *SingleThreadClient) SendRequest(req ClientRequest, address string) {
	data := req.encode(codegen.NewEncoder())
	response, err := send_request(address, data)
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
	return client.rand_int_generator.Int() % (TIMELINE_INT_MAX - TIMELINE_INT_MIN) + TIMELINE_INT_MIN
}

func (client *SingleThreadClient) _gen_text(text_len int) string {
	text := ""
	for i := 0; i < text_len; i++ {
		text += string(CHAR_SET[client.rand_charset_generator.Int() % len(CHAR_SET)])
	}
	return text
}


type HomeTimelineRequest struct {
	user_id int64
	start int
	stop int
}

func (htr *HomeTimelineRequest) encode(enc *codegen.Encoder) []byte {
	enc.Int64(htr.user_id)
	enc.Int(htr.start)
	enc.Int(htr.stop)
	return enc.Data()
}

func (htr *HomeTimelineRequest) generate(client *SingleThreadClient) {
	htr.user_id = client._gen_user_id()
	htr.start = client._gen_timeline_int()
	htr.stop = client._gen_timeline_int()
}

type UserTimelineRequest struct {
	user_id int64
	start int
	stop int
}

func (utr *UserTimelineRequest) encode(enc *codegen.Encoder) []byte {
	enc.Int64(utr.user_id)
	enc.Int(utr.start)
	enc.Int(utr.stop)
	return enc.Data()
}

func (utr *UserTimelineRequest) generate(client *SingleThreadClient) {
	utr.user_id = client._gen_user_id()
	utr.start = client._gen_timeline_int()
	utr.stop = client._gen_timeline_int()
}

type ComposePostRequest struct {
	user_id int64
	username string
	text string
	media_ids []int64
	media_types []string
	post_type PostType
}

func (cpr *ComposePostRequest) encode(enc *codegen.Encoder) []byte {
	enc.Int64(cpr.user_id)
	enc.String(cpr.username)
	enc.String(cpr.text)
	Encode_slice_int64(enc, cpr.media_ids)
	Encode_slice_string(enc, cpr.media_types)
	enc.Int((int)(cpr.post_type))
	return enc.Data()
}

func (cpr *ComposePostRequest) generate(client *SingleThreadClient) {
	cpr.user_id = client._gen_user_id()
	cpr.username = fmt.Sprintf("username_%d", cpr.user_id)
	cpr.text = client._gen_text(TEXT_LEN)

	num_mentions := client.rand_int_generator.Int() % NUM_MENTIONS_MAX
	for i := 0; i < num_mentions; i++ {
		mention_id := client._gen_user_id()
		cpr.text += fmt.Sprintf(" @username_%d", mention_id)
	}

	num_urls := client.rand_max_urls_generator.Int() % NUM_URLS_MAX
	for i := 0; i < num_urls; i++ {
		cpr.text += "http://" + client._gen_text(URL_LEN)
	}

	num_medias := client.rand_max_medias_generator.Int() % NUM_MEDIAS_MAX
	cpr.media_ids = make([]int64, num_medias)
	for i := 0; i < num_medias; i++ {
		cpr.media_ids[i] = client.rand_user_id_generator.Int63()
		cpr.media_types = append(cpr.media_types, "png")
	}
	cpr.post_type = POST
}

type RemovePostsRequest struct {
	user_id int64
	start int
	stop int
}

func (rpr *RemovePostsRequest) encode(enc *codegen.Encoder) []byte {
	enc.Int64(rpr.user_id)
	enc.Int(rpr.start)
	enc.Int(rpr.stop)
	return enc.Data()
}

func (rpr *RemovePostsRequest) generate(client *SingleThreadClient) {
	rpr.user_id = client._gen_user_id()
	rpr.start = 0
	rpr.stop = 1
}

type FollowRequest struct {
	user_id int64
	followee_id int64
}

func (fr *FollowRequest) encode(enc *codegen.Encoder) []byte {
	enc.Int64(fr.user_id)
	enc.Int64(fr.followee_id)
	return enc.Data()
}

func (fr *FollowRequest) generate(client *SingleThreadClient) {
	fr.user_id = client._gen_user_id()
	fr.followee_id = client._gen_user_id()
}

type UnfollowRequest struct {
	user_id int64
	followee_id int64
}

func (ufr *UnfollowRequest) encode(enc *codegen.Encoder) []byte {
	enc.Int64(ufr.user_id)
	enc.Int64(ufr.followee_id)
	return enc.Data()
}

func (ufr *UnfollowRequest) generate(client *SingleThreadClient) {
	ufr.user_id = client._gen_user_id()
	ufr.followee_id = client._gen_user_id()
}

type GetFollowersRequest struct {
	user_id int64
}

func (gfr *GetFollowersRequest) encode(enc *codegen.Encoder) []byte {
	enc.Int64(gfr.user_id)
	return enc.Data()
}

func (gfr *GetFollowersRequest) generate(client *SingleThreadClient) {
	gfr.user_id = client._gen_user_id()
}


func main() {
	client := SingleThreadClient{}
	client.Init()

	for {
		req, address := client.GenRequest()
		client.SendRequest(req, address)
		time.Sleep(INTERVAL_BETWEEN_REQUESTS)
	}
}