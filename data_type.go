package main

type Creator struct {
  userId int64
  username string
}

type UserProfile struct {
  userId int64
  firstName string
  lastName string
  salt string
  passwordHashed string
}
