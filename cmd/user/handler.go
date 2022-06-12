package main

import (
	"context"
	"douyin-demo-micro/cmd/user/service"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/util"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// CreateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	resp = new(user.CreateUserResponse)
	if !util.LoginCheck(req.Req.Username, req.Req.Password) {
		resp.Resp = util.FailResp(1, util.ErrIllegalArguments)
		return resp, nil
	}

	err = service.CreateUser(ctx, req)
	if err != nil {
		resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp = util.SuccessResp()
	return resp, nil
}

// CheckUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CheckUser(ctx context.Context, req *user.CheckUserRequest) (resp *user.CheckUserResponse, err error) {
	resp = new(user.CheckUserResponse)
	resp.Resp = new(user.UserResponse)

	if !util.LoginCheck(req.Req.Username, req.Req.Password) {
		resp.Resp.Resp = util.FailResp(1, util.ErrIllegalArguments)
		return resp, nil
	}
	user, err := service.CheckUser(ctx, req)
	if err != nil {
		resp.Resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp.User = user
	resp.Resp.Resp = util.SuccessResp()
	return resp, nil
}

// FindUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) FindUser(ctx context.Context, req *user.FindUserRequest) (resp *user.FindUserResponse, err error) {
	resp = new(user.FindUserResponse)
	resp.Resp = new(user.UserResponse)

	user, err := service.FindUser(ctx, req)
	if err != nil {
		resp.Resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp.User = user
	resp.Resp.Resp = util.SuccessResp()
	return resp, nil
}

// MGetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) MGetUser(ctx context.Context, req *user.MGetUserRequest) (resp *user.MGetUserResponse, err error) {
	resp = new(user.MGetUserResponse)
	resp.Resp = new(user.MultiUserResponse)

	users, err := service.MGetUser(ctx, req)
	if err != nil {
		resp.Resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp.Users = users
	resp.Resp.Resp = util.SuccessResp()
	return resp, nil
}

// Follow implements the UserServiceImpl interface.
func (s *UserServiceImpl) Follow(ctx context.Context, req *user.FollowRequest) (resp *user.FollowResponse, err error) {
	resp = new(user.FollowResponse)

	err = service.Follow(ctx, req)
	if err != nil {
		resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp = util.SuccessResp()
	return resp, nil
}

// Unfollow implements the UserServiceImpl interface.
func (s *UserServiceImpl) Unfollow(ctx context.Context, req *user.UnfollowRequest) (resp *user.UnfollowResponse, err error) {
	resp = new(user.UnfollowResponse)

	err = service.Unfollow(ctx, req)
	if err != nil {
		resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp = util.SuccessResp()
	return resp, nil
}

// FollowList implements the UserServiceImpl interface.
func (s *UserServiceImpl) FollowList(ctx context.Context, req *user.FollowListRequest) (resp *user.FollowListResponse, err error) {
	resp = new(user.FollowListResponse)
	resp.Resp = new(user.MultiUserResponse)

	users, err := service.FollowList(ctx, req)
	if err != nil {
		resp.Resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp.Users = users
	resp.Resp.Resp = util.SuccessResp()
	return resp, nil
}

// FollowerList implements the UserServiceImpl interface.
func (s *UserServiceImpl) FollowerList(ctx context.Context, req *user.FollowerListRequest) (resp *user.FollowerListResponse, err error) {
	resp = new(user.FollowerListResponse)
	resp.Resp = new(user.MultiUserResponse)

	users, err := service.FollowerList(ctx, req)
	if err != nil {
		resp.Resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp.Users = users
	resp.Resp.Resp = util.SuccessResp()
	return resp, nil
}
