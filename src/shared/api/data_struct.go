package api

import (
	"net/http"

	"SocialNetwork/shared/common"

	"github.com/ServiceWeaver/weaver/runtime/codegen"
)

type EncodableRequest interface {
	Encode(*codegen.Encoder) []byte
}

type DecodableResponse interface {
	Decode(*http.Response) []byte
}

type ClientRequest interface {
	EncodableRequest
}

type ClientResponse interface {
	DecodableResponse
}

type RegisterUserRequest struct {
	FirstName string
	LastName  string
	Username  string
	Password  string
}

func (req *RegisterUserRequest) Encode(enc *codegen.Encoder) []byte {
	enc.String(req.FirstName)
	enc.String(req.LastName)
	enc.String(req.Username)
	enc.String(req.Password)
	return enc.Data()
}

type RegisterUserWithIdRequest struct {
	FirstName string
	LastName  string
	Username  string
	Password  string
	UserId    int64
}

func (req *RegisterUserWithIdRequest) Encode(enc *codegen.Encoder) []byte {
	enc.String(req.FirstName)
	enc.String(req.LastName)
	enc.String(req.Username)
	enc.String(req.Password)
	enc.Int64(req.UserId)
	return enc.Data()
}

type ReadHomeTimelineRequest struct {
	UserId int64
	Start  int
	Stop   int
}

func (rhtr *ReadHomeTimelineRequest) Encode(enc *codegen.Encoder) []byte {
	enc.Int64(rhtr.UserId)
	enc.Int(rhtr.Start)
	enc.Int(rhtr.Stop)
	return enc.Data()
}

type ReadUserTimelineRequest struct {
	UserId int64
	Start  int
	Stop   int
}

func (rutr *ReadUserTimelineRequest) Encode(enc *codegen.Encoder) []byte {
	enc.Int64(rutr.UserId)
	enc.Int(rutr.Start)
	enc.Int(rutr.Stop)
	return enc.Data()
}

type ComposePostRequest struct {
	Username   string
	UserId     int64
	Text       string
	MediaIds   []int64
	MediaTypes []string
	PostType   common.PostType
}

func (cpr *ComposePostRequest) Encode(enc *codegen.Encoder) []byte {
	enc.String(cpr.Username)
	enc.Int64(cpr.UserId)
	enc.String(cpr.Text)
	common.Encode_slice_int64(enc, cpr.MediaIds)
	common.Encode_slice_string(enc, cpr.MediaTypes)
	enc.Int((int)(cpr.PostType))
	return enc.Data()
}

type RemovePostsRequest struct {
	UserId int64
	Start  int
	Stop   int
}

func (rpr *RemovePostsRequest) Encode(enc *codegen.Encoder) []byte {
	enc.Int64(rpr.UserId)
	enc.Int(rpr.Start)
	enc.Int(rpr.Stop)
	return enc.Data()
}

type LoginRequest struct {
	Username string
	Password string
}

func (req *LoginRequest) Encode(enc *codegen.Encoder) []byte {
	enc.String(req.Username)
	enc.String(req.Password)
	return enc.Data()
}

type FollowRequest struct {
	UserId     int64
	FolloweeId int64
}

func (fr *FollowRequest) Encode(enc *codegen.Encoder) []byte {
	enc.Int64(fr.UserId)
	enc.Int64(fr.FolloweeId)
	return enc.Data()
}

type FollowWithUsernameRequest struct {
	Username         string
	FolloweeUsername string
}

func (req *FollowWithUsernameRequest) Encode(enc *codegen.Encoder) []byte {
	enc.String(req.Username)
	enc.String(req.FolloweeUsername)
	return enc.Data()
}

type UnfollowRequest struct {
	UserId     int64
	FolloweeId int64
}

func (ufr *UnfollowRequest) Encode(enc *codegen.Encoder) []byte {
	enc.Int64(ufr.UserId)
	enc.Int64(ufr.FolloweeId)
	return enc.Data()
}

type UnfollowWithUsernameRequest struct {
	Username        string
	FolloweeUsernae string
}

func (req *UnfollowWithUsernameRequest) Encode(enc *codegen.Encoder) []byte {
	enc.String(req.Username)
	enc.String(req.FolloweeUsernae)
	return enc.Data()
}

type GetFollowersRequest struct {
	UserId int64
}

func (gfr *GetFollowersRequest) Encode(enc *codegen.Encoder) []byte {
	enc.Int64(gfr.UserId)
	return enc.Data()
}

type GetFolloweesRequest struct {
	UserId int64
}

func (req *GetFolloweesRequest) Encode(enc *codegen.Encoder) []byte {
	enc.Int64(req.UserId)
	return enc.Data()
}

type UploadMediaRequest struct {
	Filename string
	Data     string
}

func (req *UploadMediaRequest) Encode(enc *codegen.Encoder) []byte {
	enc.String(req.Filename)
	enc.String(req.Data)
	return enc.Data()
}

type GetMediaRequest struct {
	Filename string
}

func (req *GetMediaRequest) Encode(enc *codegen.Encoder) []byte {
	enc.String(req.Filename)
	return enc.Data()
}
