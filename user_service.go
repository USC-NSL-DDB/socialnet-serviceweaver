package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/golang-jwt/jwt"
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
    ComposeCreatorWithUsername(context.Context, string) Creator
    ComposeCreatorWithUserId(context.Context, int64, string) Creator

    Login(context.Context, string, string) (string, error)
    GetUserId(context.Context, string) int64
    
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

func (us *UserService) Init(context.Context) error {
  us.LoadSecretAndMachineId()
  return nil 
}

func (us *UserService) LoadSecretAndMachineId() {
  // figure out how to load data from local config file
}

func (us *UserService) ComposeCreatorWithUsername(ctx context.Context, username string) Creator {
  storage := us.storage.Get()
  profile, exist := storage.GetUserProfile(ctx, username)
  if !exist {
    fmt.Printf("Failed to find the user profile - username: %s\n", username)
    // Should handle it as error
    return Creator{
      userId: 0,
      username: "",
    }
  }

  return us.ComposeCreatorWithUserId(ctx, profile.userId, username)
}

func (us *UserService) ComposeCreatorWithUserId(ctx context.Context, userId int64, username string) Creator {
  return Creator {
    userId: userId,
    username: username,
  }
}

func (us *UserService) Login(ctx context.Context, username, password string) (string, error) {
  storage := us.storage.Get()
  profile, exist := storage.GetUserProfile(ctx, username)
  if !exist {
    return "", &LogError{
      err_code: int(NOT_REGISTERED),
      err_msg: NOT_REGISTERED.String(),
    }
  }
  var auth bool = HashPassowrd(password, profile.salt) == profile.passwordHashed
  if !auth {
    return "", &LogError{
      err_code: int(WRONG_PASSWORD),
      err_msg: WRONG_PASSWORD.String(),
    }
  }

  userIdStr := strconv.FormatInt(profile.userId, 10)
	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)

	// Create a new JWT object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userIdStr,
		"username":  username,
		"timestamp": timestampStr,
		"ttl":       "3600",
	})

  secret := "mysecret"

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}

  return tokenString, nil
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

func (us *UserService) GetUserId(ctx context.Context, username string) int64 {
  storage := us.storage.Get()
  profile, exist := storage.GetUserProfile(ctx, username)
  if !exist {
    fmt.Printf("Err. no profile associated with username: %s.\n", username)
    // Should handle it more elegantly.
    return 0 
  }
  return profile.userId
}


