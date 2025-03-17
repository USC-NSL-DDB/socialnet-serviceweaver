package main

import (
	"SocialNetwork/shared/api"
	"bytes"
	"encoding/binary"
)

// Custom encoder to match your service's encoding
type Encoder struct {
	buf bytes.Buffer
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) String(s string) {
	binary.Write(&e.buf, binary.LittleEndian, int32(len(s)))
	e.buf.WriteString(s)
}

func main() {
	// Configuration
	serviceURL := "http://localhost:8081" // Adjust this to your service's address
	api.RegisterUser(serviceURL, &api.RegisterUserRequest{
		Username:  "username_74",
		Password:  "password_74",
		FirstName: "first_name_74",
		LastName:  "last_name_74",
	})
	api.Login(serviceURL, &api.LoginRequest{
		Username: "username_74",
		Password: "password_74",
	})

}
