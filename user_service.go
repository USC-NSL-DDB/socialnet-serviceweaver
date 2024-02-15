package main

import (
	"context"
  "time"
	"fmt"
  "math/rand"
  "crypto/sha256"
  "encoding/hex"

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
    RegisterUserWithId(context.Context, string, string, string, string, int64) 
    RegisterUser(context.Context, string, string, string, string) 
    // TODO: Figure out what is Creator return type
    // ComposeCreatorWithUsername(string) Creator
    // ComposeCreatorWithUserId(int64, string) Creator

    Login(context.Context, string, string) (string, error)
    
}

func GenRandomString(length int) string {
  const alphanum string = `0123456789
                           ABCDEFGHIJKLMNOPQRSTUVWXYZ
                           abcdefghijklmnopqrstuvwxyz`
  s := ""
  for i := 0; i < length; i++ {
    rd_idx := rand.Intn(len(alphanum))
    s += string(alphanum[rd_idx])
  }
  return s
}

func HashPassowrd(password, salt string) string {
  combined := salt

  hasher := sha256.New()
  hasher.Write([]byte(combined))
  hashBytes := hasher.Sum(nil)

  hashHex := hex.EncodeToString(hashBytes)
  return hashHex
}

func GenerateUniqueId() int64 {
    // Get the current Unix timestamp in milliseconds
    // This reduces the chance of collision for IDs generated in quick succession
    timestamp := time.Now().UnixNano() / int64(time.Millisecond)
    
    // Shift the timestamp to occupy the higher bits of the int64, making room for the random component
    // Adjust the shifting based on your application's needs for timestamp precision vs. randomness
    timestamp = timestamp << (64 - 48) // Adjust 48 based on your needs
    
    // Generate a random component to fill the lower bits
    // Ensure the random source is seeded (usually done once globally)
    randComponent := rand.Int63n(1 << 48) // Adjust 48 based on the shifting above
    
    // Combine the timestamp and random component
    uniqueID := timestamp | randComponent
    
    return uniqueID
}

type UserService struct {
    weaver.Implements[UserServicer]
    storage weaver.Ref[Storage]

    _machineId string
    _secret string
}

// func (us *UserService) Init(context.Context) error {
//   
// }

func (us *UserService) LoadSecretAndMachineId() {
  // figure out how to load data from local config file
}

func (us *UserService) Login(_ context.Context, username, password string) (string, error) {
  return "", nil
}

func (us *UserService) RegisterUserWithId(ctx context.Context, firstName, lastName, username, password string, userId int64) {
  salt := GenRandomString(32)
  userProfile := UserProfile {
    userId: userId,
    firstName: firstName,
    lastName: lastName,
    salt: salt,
    passwordHashed: HashPassowrd(password, salt),
  }
  // update the map
  var s Storage = us.storage.Get()
  s.PutUserProfile(ctx, username, userProfile)
}

func (us *UserService) RegisterUser(ctx context.Context, firstName, lastName, username, password string) {
  // Generate a user id
  uid := GenerateUniqueId()

  // Call RegisterUserWithId
  us.RegisterUserWithId(ctx, firstName, lastName, username, password, uid)
}

type UserProfile struct {
  userId int64
  firstName string
  lastName string
  salt string
  passwordHashed string
}

