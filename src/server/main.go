package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	. "SocialNetwork/shared/common"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
)

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	backend_service weaver.Ref[BackendServicer]

	remove_posts           weaver.Listener
	compose_post           weaver.Listener
	login                  weaver.Listener
	register_user          weaver.Listener
	register_user_with_id  weaver.Listener
	read_user_timeline     weaver.Listener
	get_followers          weaver.Listener
	unfollow               weaver.Listener
	unfollow_with_username weaver.Listener
	follow                 weaver.Listener
	follow_with_username   weaver.Listener
	get_followees          weaver.Listener
	read_home_timeline     weaver.Listener
	upload_media           weaver.Listener
	get_media              weaver.Listener
}

func reg_listener_action(
	listener weaver.Listener,
	endpoint string,
	action func(http.ResponseWriter, *http.Request),
	error_collector chan error,
) {
	go func() {
		fmt.Printf("%v available on %v\n", endpoint, listener)
		http.HandleFunc(endpoint, action)
		err := http.Serve(listener, nil)
		error_collector <- err
	}()
}

func decode_request_body(r *http.Request, action func(*codegen.Decoder)) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	dec := codegen.NewDecoder(body)
	action(dec)
	return nil
}

func encode_response_body(w http.ResponseWriter, action func(*codegen.Encoder)) {
	w.Header().Set("Content-Type", "application/custom")
	enc := codegen.NewEncoder()
	action(enc)
	w.Write(enc.Data())
}

// serve is called by weaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {
	var backend = app.backend_service.Get()
	err_collector := make(chan error)

	reg_listener_action(app.remove_posts, "/remove_posts", func(w http.ResponseWriter, r *http.Request) {
		var user_id int64
		var start int
		var stop int

		decode_request_body(r, func(dec *codegen.Decoder) {
			user_id = dec.Int64()
			start = dec.Int()
			stop = dec.Int()
		})

		err := backend.RemovePosts(context.Background(), user_id, start, stop)
		if err != nil {
			log.Default().Println(err)
		}
		// fmt.Fprintf(w, "remove_posts\n")
	}, err_collector)

	reg_listener_action(app.compose_post, "/compose_post", func(w http.ResponseWriter, r *http.Request) {
		var username string
		var user_id int64
		var text string
		var media_ids []int64
		var media_types []string
		var post_type PostType

		decode_request_body(r, func(dec *codegen.Decoder) {
			username = dec.String()
			user_id = dec.Int64()
			text = dec.String()
			media_ids = Decode_slice_int64(dec)
			media_types = Decode_slice_string(dec)
			post_type = (PostType)(dec.Int())
		})

		err := backend.CompostPost(
			context.Background(),
			username, user_id, text, media_ids, media_types, post_type,
		)
		if err != nil {
			log.Default().Println(err)
			// r.Response.StatusCode = 500
		}
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.login, "/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.register_user, "/register_user", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.register_user_with_id, "/register_user_with_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.read_user_timeline, "/read_user_timeline", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.get_followers, "/get_followers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.unfollow, "/unfollow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.unfollow_with_username, "/unfollow_with_username", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.follow, "/follow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.follow_with_username, "/follow_with_username", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.get_followees, "/get_followees", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.read_home_timeline, "/read_home_timeline", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.upload_media, "/upload_media", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	reg_listener_action(app.get_media, "/get_media", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "compose_post\n")
	}, err_collector)

	for err := range err_collector {
		log.Fatal(err)
		return err
	}

	return nil
}
