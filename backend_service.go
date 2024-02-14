package main

import (
    "context"

    "github.com/ServiceWeaver/weaver"
)


type BackendServicer interface {

    RemovePost(context.Context, int64, int)
    CompostPost(context.Context, string, int64, string, []int64, []string)
    Login(context.Context, string, string, string)
    // RegisterUser(context.Context)
    // ReadUserTimeline(context.Context, )
    // Reverse(context.Context, string) (string, error)
}

type BackendService struct{
    weaver.Implements[BackendServicer]
}

// func (r *reverser) Reverse(_ context.Context, s string) (string, error) {
//     runes := []rune(s)
//     n := len(runes)
//     for i := 0; i < n/2; i++ {
//         runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
//     }
//     return string(runes), nil
// }
