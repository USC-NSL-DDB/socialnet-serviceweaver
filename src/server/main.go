package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"SocialNetwork/shared/common"

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

	// remove_posts           weaver.Listener
	// compose_post           weaver.Listener
	// login                  weaver.Listener
	// register_user          weaver.Listener
	// register_user_with_id  weaver.Listener
	// read_user_timeline     weaver.Listener
	// get_followers          weaver.Listener
	// unfollow               weaver.Listener
	// unfollow_with_username weaver.Listener
	// follow                 weaver.Listener
	// follow_with_username   weaver.Listener
	// get_followees          weaver.Listener
	// read_home_timeline     weaver.Listener
	// upload_media           weaver.Listener
	// get_media              weaver.Listener
	api_listener weaver.Listener
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

	reg_listener_action(app.api_listener, common.REMOVE_POSTS_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
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

	reg_listener_action(app.api_listener, common.COMPOSE_POST_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
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
			media_ids = common.Decode_slice_int64(dec)
			media_types = common.Decode_slice_string(dec)
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

	reg_listener_action(app.api_listener, common.LOGIN_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var username string
		var password string

		decode_request_body(r, func(dec *codegen.Decoder) {
			username = dec.String()
			password = dec.String()
		})

		token, err := backend.Login(context.Background(), username, password)
		if err != nil {
			log.Default().Println(err) // never triggered
		} else {
			encode_response_body(w, func(enc *codegen.Encoder) {
				enc.String(token)
			})
		}

		fmt.Fprintf(w, "login\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.REGISTER_USER_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var first_name string
		var last_name string
		var username string
		var password string

		decode_request_body(r, func(dec *codegen.Decoder) {
			first_name = dec.String()
			last_name = dec.String()
			username = dec.String()
			password = dec.String()
		})

		err := backend.RegisterUser(context.Background(), first_name, last_name, username, password)
		if err != nil {
			log.Default().Println(err) // never triggered
		}

		fmt.Fprintf(w, "register_user\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.REGISTER_USER_WITH_ID_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var first_name string
		var last_name string
		var username string
		var password string
		var user_id int64

		decode_request_body(r, func(dec *codegen.Decoder) {
			first_name = dec.String()
			last_name = dec.String()
			username = dec.String()
			password = dec.String()
			user_id = dec.Int64()
		})

		err := backend.RegisterUserWithId(context.Background(), first_name, last_name, username, password, user_id)
		if err != nil {
			log.Default().Println(err) // never triggered
		}

		fmt.Fprintf(w, "register_user_with_id\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.READ_USER_TIMELINE_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var user_id int64
		var start int
		var stop int

		decode_request_body(r, func(dec *codegen.Decoder) {
			user_id = dec.Int64()
			start = dec.Int()
			stop = dec.Int()
		})

		posts, err := backend.ReadUserTimeline(context.Background(), user_id, start, stop)
		if err != nil {
			log.Default().Println(err)
		} else {
			encode_response_body(w, func(enc *codegen.Encoder) {
				enc.Int(len(posts)) // TODO: repeated code in read_home_timeline
				for _, post := range posts {
					enc.Int64(post.Post_id)
					enc.Int64(post.Creator.UserId)
					enc.String(post.Creator.Username)
					enc.Int64(post.Req_id)
					enc.String(post.Text)
					enc.Int64(post.Timestamp)
					enc.Int(int(post.Post_type))

					enc.Int(len(post.User_mentions))
					enc.Int(len(post.Media))
					enc.Int(len(post.Urls))
					for _, user_mention := range post.User_mentions {
						enc.Int64(user_mention.UserId)
						enc.String(user_mention.Username)
					}
					for _, media := range post.Media {
						enc.Int64(media.MediaId)
						enc.String(media.MediaType)
					}
					for _, url := range post.Urls {
						enc.String(url.ShortenedUrl) // send only shortened url, check if it is correct
					}
				}
			})
		}

		fmt.Fprintf(w, "read_user_timeline\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.GET_FOLLOWERS_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var user_id int64

		decode_request_body(r, func(dec *codegen.Decoder) {
			user_id = dec.Int64()
		})

		followers, err := backend.GetFollowers(context.Background(), user_id)
		if err != nil {
			log.Default().Println(err)
		} else {
			encode_response_body(w, func(enc *codegen.Encoder) {
				enc.Int(len(followers))
				for _, follower_id := range followers {
					enc.Int64(follower_id)
				}
			})
		}

		fmt.Fprintf(w, "get_followers\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.UNFOLLOW_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var id int64
		var followee_id int64

		decode_request_body(r, func(dec *codegen.Decoder) {
			id = dec.Int64()
			followee_id = dec.Int64()
		})

		err := backend.Unfollow(context.Background(), id, followee_id)
		if err != nil {
			log.Default().Println(err) // never triggered
		}

		fmt.Fprintf(w, "unfollow\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.UNFOLLOW_WITH_USERNAME_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var username string
		var followee_username string

		decode_request_body(r, func(dec *codegen.Decoder) {
			username = dec.String()
			followee_username = dec.String()
		})

		err := backend.UnfollowWithUsername(context.Background(), username, followee_username)
		if err != nil {
			log.Default().Println(err) // never triggered
		}

		fmt.Fprintf(w, "unfollow_with_username\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.FOLLOW_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var id int64
		var followee_id int64

		decode_request_body(r, func(dec *codegen.Decoder) {
			id = dec.Int64()
			followee_id = dec.Int64()
		})

		err := backend.Follow(context.Background(), id, followee_id)
		if err != nil {
			log.Default().Println(err) // never triggered
		}

		fmt.Fprintf(w, "follow\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.FOLLOW_WITH_USERNAME_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var username string
		var followee_username string

		decode_request_body(r, func(dec *codegen.Decoder) {
			username = dec.String()
			followee_username = dec.String()
		})

		err := backend.FollowWithUsername(context.Background(), username, followee_username)
		if err != nil {
			log.Default().Println(err) // never triggered
		}

		fmt.Fprintf(w, "follow_with_username\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.GET_FOLLOWEES_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var user_id int64

		decode_request_body(r, func(dec *codegen.Decoder) {
			user_id = dec.Int64()
		})

		followees, err := backend.GetFollowees(context.Background(), user_id)
		if err != nil {
			log.Default().Println(err)
		} else {
			encode_response_body(w, func(enc *codegen.Encoder) {
				enc.Int(len(followees))
				for _, followee_id := range followees {
					enc.Int64(followee_id)
				}
			})
		}

		fmt.Fprintf(w, "get_followees\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.READ_HOME_TIMELINE_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var user_id int64
		var start int
		var stop int

		decode_request_body(r, func(dec *codegen.Decoder) {
			user_id = dec.Int64()
			start = dec.Int()
			stop = dec.Int()
		})

		posts, err := backend.ReadHomeTimeline(context.Background(), user_id, start, stop)
		if err != nil {
			log.Default().Println(err)
		} else {
			encode_response_body(w, func(enc *codegen.Encoder) {
				enc.Int(len(posts)) // TODO: repeated code in read_user_timeline
				for _, post := range posts {
					enc.Int64(post.Post_id)
					enc.Int64(post.Creator.UserId)
					enc.String(post.Creator.Username)
					enc.Int64(post.Req_id)
					enc.String(post.Text)
					enc.Int64(post.Timestamp)
					enc.Int(int(post.Post_type))

					enc.Int(len(post.User_mentions))
					enc.Int(len(post.Media))
					enc.Int(len(post.Urls))
					for _, user_mention := range post.User_mentions {
						enc.Int64(user_mention.UserId)
						enc.String(user_mention.Username)
					}
					for _, media := range post.Media {
						enc.Int64(media.MediaId)
						enc.String(media.MediaType)
					}
					for _, url := range post.Urls {
						enc.String(url.ShortenedUrl) // send only shortened url, check if it is correct
					}
				}
			})
		}

		fmt.Fprintf(w, "read_home_timeline\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.UPLOAD_MEDIA_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var filename string
		var data string

		decode_request_body(r, func(dec *codegen.Decoder) {
			filename = dec.String()
			data = dec.String()
		})

		err := backend.UploadMedia(context.Background(), filename, data)
		if err != nil {
			log.Default().Println(err) // never triggered
		}
		fmt.Fprintf(w, "upload_media\n")
	}, err_collector)

	reg_listener_action(app.api_listener, common.GET_MEDIA_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		var filename string

		decode_request_body(r, func(dec *codegen.Decoder) {
			filename = dec.String()
		})

		media, err := backend.GetMedia(context.Background(), filename)
		if err != nil {
			log.Default().Println(err)
		} else {
			encode_response_body(w, func(enc *codegen.Encoder) {
				enc.String(media)
			})
		}

		fmt.Fprintf(w, "get_media\n")
	}, err_collector)

	for err := range err_collector {
		log.Fatal(err)
		return err
	}

	return nil
}
