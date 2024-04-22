package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"SocialNetwork/shared/common"

	"github.com/ServiceWeaver/weaver/runtime/codegen"
)

const (
	BASE_PORT = "49555"
	BASE_URL  = "http://localhost:" + BASE_PORT
)

func send_request(address string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/custom")
	client := &http.Client{}
	return client.Do(req)
}

func InitRequest(endpoint string, data []byte) (*http.Response, error) {
	address := BASE_URL + endpoint
	return send_request(address, data)
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

func RemovePosts(req *RemovePostsRequest) {
	resp, err := InitRequest(common.REMOVE_POSTS_ENDPOINT, req.Encode(codegen.NewEncoder()))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}

func ComposePost(req *ComposePostRequest) {
	resp, err := InitRequest(common.COMPOSE_POST_ENDPOINT, req.Encode(codegen.NewEncoder()))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}
