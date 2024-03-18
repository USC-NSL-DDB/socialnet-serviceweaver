package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
)

func main() {
	if err := weaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

// app is the main component of the application. weaver.Run creates
// it and passes it to serve.
type app struct {
	weaver.Implements[weaver.Main]
	// reverser weaver.Ref[Reverser]
	backend_service weaver.Ref[BackendServicer]

	remove_posts weaver.Listener
	compose_post weaver.Listener
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
}

func reg_listener_action(
	listener weaver.Listener,
	endpoint string,
	action func(http.ResponseWriter, *http.Request),
) {
	go func() {
		fmt.Printf("%v available on %v\n", endpoint, listener)
		http.HandleFunc(endpoint, action)
		err := http.Serve(listener, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	// fmt.Printf("%v available on %v\n", endpoint, listener)
	// http.HandleFunc(endpoint, action)
	// err := http.Serve(listener, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// serve is called by weaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {
	// var backend BackendService = app.backend_service.Get()

	reg_listener_action(app.remove_posts, "/remove_posts", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "remove_posts")
	})

	// reg_listener_action(app.compose_post, "/compose_post", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "compose_post")
	// })
	// var r Reverser = app.reverser.Get()
	// reversed, err := r.Reverse(ctx, "!dlroW ,olleH")
	// if err != nil {
	//   return err
	// }
	// fmt.Println(reversed)
	// return nil
	return nil
}
