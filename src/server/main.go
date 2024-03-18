package main

import (
	"context"
	"log"

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
	backend_service weaver.Ref[BackendService]
}

// serve is called by weaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {
	var backend BackendService = app.backend_service.Get()
	backend.Follow(ctx, 1, 2)
	// var r Reverser = app.reverser.Get()
	// reversed, err := r.Reverse(ctx, "!dlroW ,olleH")
	// if err != nil {
	//   return err
	// }
	// fmt.Println(reversed)
	// return nil
	return nil
}
