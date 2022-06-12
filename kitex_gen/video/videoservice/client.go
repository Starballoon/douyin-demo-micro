// Code generated by Kitex v0.3.1. DO NOT EDIT.

package videoservice

import (
	"context"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/kitex_gen/video"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	CreateVideo(ctx context.Context, req *video.CreateVideoRequest, callOptions ...callopt.Option) (r *user.BaseResp, err error)
	Feed(ctx context.Context, req *video.FeedRequest, callOptions ...callopt.Option) (r *video.FeedResponse, err error)
	MGetVideo(ctx context.Context, red *video.MGetVideoRequest, callOptions ...callopt.Option) (r *video.MGetVideoResponse, err error)
	PublishList(ctx context.Context, req *video.PublishListRequest, callOptions ...callopt.Option) (r *video.PublishListResponse, err error)
	FavoriteList(ctx context.Context, req *video.FavoriteListRequest, callOptions ...callopt.Option) (r *video.FavoriteListResponse, err error)
	Favorite(ctx context.Context, req *video.FavoriteRequest, callOptions ...callopt.Option) (r *video.FavoriteResponse, err error)
	Unfavorite(ctx context.Context, req *video.UnfavoriteRequest, callOptions ...callopt.Option) (r *video.UnfavoriteResponse, err error)
	CheckVideo(ctx context.Context, req *video.CheckVideoRequest, callOptions ...callopt.Option) (r *user.BaseResp, err error)
	UpdateCommentCount(ctx context.Context, req *video.UpdateCommentCountRequest, callOptions ...callopt.Option) (r *user.BaseResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kVideoServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kVideoServiceClient struct {
	*kClient
}

func (p *kVideoServiceClient) CreateVideo(ctx context.Context, req *video.CreateVideoRequest, callOptions ...callopt.Option) (r *user.BaseResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateVideo(ctx, req)
}

func (p *kVideoServiceClient) Feed(ctx context.Context, req *video.FeedRequest, callOptions ...callopt.Option) (r *video.FeedResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Feed(ctx, req)
}

func (p *kVideoServiceClient) MGetVideo(ctx context.Context, red *video.MGetVideoRequest, callOptions ...callopt.Option) (r *video.MGetVideoResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MGetVideo(ctx, red)
}

func (p *kVideoServiceClient) PublishList(ctx context.Context, req *video.PublishListRequest, callOptions ...callopt.Option) (r *video.PublishListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PublishList(ctx, req)
}

func (p *kVideoServiceClient) FavoriteList(ctx context.Context, req *video.FavoriteListRequest, callOptions ...callopt.Option) (r *video.FavoriteListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.FavoriteList(ctx, req)
}

func (p *kVideoServiceClient) Favorite(ctx context.Context, req *video.FavoriteRequest, callOptions ...callopt.Option) (r *video.FavoriteResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Favorite(ctx, req)
}

func (p *kVideoServiceClient) Unfavorite(ctx context.Context, req *video.UnfavoriteRequest, callOptions ...callopt.Option) (r *video.UnfavoriteResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Unfavorite(ctx, req)
}

func (p *kVideoServiceClient) CheckVideo(ctx context.Context, req *video.CheckVideoRequest, callOptions ...callopt.Option) (r *user.BaseResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CheckVideo(ctx, req)
}

func (p *kVideoServiceClient) UpdateCommentCount(ctx context.Context, req *video.UpdateCommentCountRequest, callOptions ...callopt.Option) (r *user.BaseResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UpdateCommentCount(ctx, req)
}
