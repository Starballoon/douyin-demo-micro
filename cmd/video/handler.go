package main

import (
	"context"
	"douyin-demo-micro/cmd/video/service"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/kitex_gen/video"
	"douyin-demo-micro/util"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// CreateVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) CreateVideo(ctx context.Context, req *video.CreateVideoRequest) (resp *user.BaseResp, err error) {
	resp = new(user.BaseResp)

	err = service.CreateVideo(ctx, req)
	if err != nil {
		resp = util.FailResp(1, err)
		return resp, nil
	}
	resp = util.SuccessResp()
	return resp, nil
}

// Feed implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	resp = new(video.FeedResponse)
	resp.Resp = new(user.BaseResp)

	videos, nextTime, err := service.Feed(ctx, req)
	if err != nil {
		resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.VideoList = videos
	resp.NextTime = nextTime
	resp.Resp = util.SuccessResp()
	return resp, nil
}

// MGetVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) MGetVideo(ctx context.Context, red *video.MGetVideoRequest) (resp *video.MGetVideoResponse, err error) {
	// TODO: Your code here...
	return
}

// PublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	resp = new(video.PublishListResponse)
	resp.Resp = new(video.MultiVideoResponse)
	resp.Resp.Resp = new(user.BaseResp)

	videos, err := service.PublishList(ctx, req)
	if err != nil {
		resp.Resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp.Videos = videos
	resp.Resp.Resp = util.SuccessResp()
	return resp, nil
}

// FavoriteList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) FavoriteList(ctx context.Context, req *video.FavoriteListRequest) (resp *video.FavoriteListResponse, err error) {
	resp = new(video.FavoriteListResponse)
	resp.Resp = new(video.MultiVideoResponse)
	resp.Resp.Resp = new(user.BaseResp)

	videos, err := service.FavoriteList(ctx, req)
	if err != nil {
		resp.Resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp.Videos = videos
	resp.Resp.Resp = util.SuccessResp()
	return resp, nil
}

// Favorite implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Favorite(ctx context.Context, req *video.FavoriteRequest) (resp *video.FavoriteResponse, err error) {
	resp = new(video.FavoriteResponse)
	resp.Resp = new(user.BaseResp)

	err = service.Favorite(ctx, req)
	if err != nil {
		resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp = util.SuccessResp()
	return resp, nil
}

// Unfavorite implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) Unfavorite(ctx context.Context, req *video.UnfavoriteRequest) (resp *video.UnfavoriteResponse, err error) {
	resp = new(video.UnfavoriteResponse)
	resp.Resp = new(user.BaseResp)

	err = service.Unfavorite(ctx, req)
	if err != nil {
		resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Resp = util.SuccessResp()
	return resp, nil
}

// CheckVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) CheckVideo(ctx context.Context, req *video.CheckVideoRequest) (resp *user.BaseResp, err error) {
	resp = new(user.BaseResp)

	err = service.CheckVideo(ctx, req)
	if err != nil {
		resp = util.FailResp(1, err)
		return resp, nil
	}
	resp = util.SuccessResp()
	return resp, nil
}

// UpdateCommentCount implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) UpdateCommentCount(ctx context.Context, req *video.UpdateCommentCountRequest) (resp *user.BaseResp, err error) {
	resp = new(user.BaseResp)

	err = service.UpdateCommentCount(ctx, req)
	if err != nil {
		resp = util.FailResp(1, err)
		return resp, nil
	}
	resp = util.SuccessResp()
	return resp, nil
}
