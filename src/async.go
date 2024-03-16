package main

import (
	"context"
)

type IFuture interface {
	Await() interface{}
}

type Future struct {
	await func(ctx context.Context) interface{}
}

func (f Future) Await() interface{} {
	return f.await(context.Background())
}

// Exec executes the async function
func AsyncExec(f func() interface{}) IFuture {
	var result interface{}
	c := make(chan struct{})
	go func() {
		defer close(c)
		result = f()
	}()
	return Future{
		await: func(ctx context.Context) interface{} {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return result
			}
		},
	}
}
