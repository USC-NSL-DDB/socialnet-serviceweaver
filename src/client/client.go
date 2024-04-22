package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	. "SocialNetwork/shared/common"

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

func init_request(endpoint string, data []byte) (*http.Response, error) {
	address := BASE_URL + endpoint
	return send_request(address, data)
}

func encode_data(action func(*codegen.Encoder)) []byte {
	enc := codegen.NewEncoder()
	action(enc)
	return enc.Data()
}

func decode_data(response *http.Response, action func(*codegen.Decoder)) {
	defer response.Body.Close()
	// fmt.Println("Status Code:", response.StatusCode)

	resp_body, _ := io.ReadAll(response.Body)
	dec := codegen.NewDecoder(resp_body)
	action(dec)
}

func send_remove_posts(user_id int64, start int, stop int) {
	data := encode_data(func(enc *codegen.Encoder) {
		enc.Int64(user_id)
		enc.Int(start)
		enc.Int(stop)
	})

	address := "http://localhost:49555" + REMOVE_POSTS_ENDPOINT

	response, err := send_request(address, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()
}

func send_compose_post(
	username string, user_id int64, text string, media_ids []int64, media_types []string, post_type PostType,
) {
	data := encode_data(func(enc *codegen.Encoder) {
		enc.String(username)
		enc.Int64(user_id)
		enc.String(text)
		Encode_slice_int64(enc, media_ids)
		Encode_slice_string(enc, media_types)
		enc.Int((int)(post_type))
	})

	address := "http://localhost:49555" + COMPOSE_POST_ENDPOINT

	response, err := send_request(address, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()
	// fmt.Println("Status Code:", response.StatusCode)
}

func main() {
	// send_compose_post()
	send_remove_posts(0, 1, 2)
}
