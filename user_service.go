package main

import (
	"context"
	"fmt"

	"github.com/ServiceWeaver/weaver"
)

type LogErrorCode int

const (
  NOT_REGISTERED LogErrorCode = iota + 1
  WRONG_PASSWORD
)

func (lec LogErrorCode) String() string {
	return [...]string{"NOT REGISTERED", "WRONG PASSWORD"}[lec-1]
}

type LogError struct {
  err_code int
  err_msg string
}

func NewLogError(err_code LogErrorCode) LogError {
  return LogError {
    err_code: int(err_code),
    err_msg: err_code.String(),
  }
}

func (le *LogError) Error() string {
  return fmt.Sprintf("Log error. code: %d, err: %v", le.err_code, le.err_msg)
}

type UserServicer interface {
    RegisterUserWithId(string, string, string, string, int64) 
    RegisterUser(string, string, string, string) 
    // TODO: Figure out what is Creator return type
    // ComposeCreatorWithUsername(string) Creator
    // ComposeCreatorWithUserId(int64, string) Creator

    Login(string, string) (string, error)
    
}

type UserService struct {
    weaver.Implements[UserServicer]

    _machine_id string
    _secret string
    _username_to_userprofile_map hashtable
}

func (us *UserService) LoadSecretAndMachineId() {
  // figure out how to load data from local config file
}

func (us *UserService) Login(username, password string) (string, error) {
  return "", nil
}

func (us *UserService) RegisterUserWithId(first_name, last_name, username, password string, user_id int64) {
  u_profile := UserProfile {
    first_name: first_name,
    last_name: last_name,
    user_id: user_id,
    // salt
    // password by hashing with salt
  }
  _ = u_profile
  // update the map
}

func (us *UserService) RegisterUser(first_name, last_name, username, password string) {
  // Generate a user id
  // Call RegisterUserWithId
}

type UserProfile struct {
  user_id int64
  first_name string
  last_name string
  salt string
  password_hased string
}


