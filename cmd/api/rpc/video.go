package rpc

import (
	"context"
	"douyin-demo-micro/kitex_gen/video"
	"douyin-demo-micro/kitex_gen/video/videoservice"
	"douyin-demo-micro/util"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"time"
)

var videoClient videoservice.Client

func initVideoRPC() {
	r, err := etcd.NewEtcdResolver([]string{util.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := videoservice.NewClient(
		util.VideoService,
		client.WithMiddleware(util.CommonMiddleware),
		client.WithInstanceMW(util.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(trace.NewDefaultClientSuite()),   // tracer
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	videoClient = c
}

func Feed(ctx context.Context, req *video.FeedRequest) ([]*video.Video, *int64, error) {
	resp, err := videoClient.Feed(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	if resp.Resp.StatusCode != 0 {
		return nil, nil, util.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMessage)
	}
	resolveEndpoint(resp.VideoList)
	return resp.VideoList, resp.NextTime, nil
}

func FavoriteList(ctx context.Context, req *video.FavoriteListRequest) ([]*video.Video, error) {
	resp, err := videoClient.FavoriteList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.Resp.StatusCode, resp.Resp.Resp.StatusMessage)
	}
	resolveEndpoint(resp.Resp.Videos)
	return resp.Resp.Videos, nil
}

func PublishList(ctx context.Context, req *video.PublishListRequest) ([]*video.Video, error) {
	resp, err := videoClient.PublishList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.Resp.StatusCode, resp.Resp.Resp.StatusMessage)
	}
	resolveEndpoint(resp.Resp.Videos)
	return resp.Resp.Videos, nil
}

func PublishVideo(ctx context.Context, req *video.CreateVideoRequest) error {
	resp, err := videoClient.CreateVideo(ctx, req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 0 {
		return util.NewErrNo(resp.StatusCode, resp.StatusMessage)
	}
	return nil
}

func Favorite(ctx context.Context, req *video.FavoriteRequest) error {
	resp, err := videoClient.Favorite(ctx, req)
	if err != nil {
		return err
	}
	if resp.Resp.StatusCode != 0 {
		return util.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMessage)
	}
	return nil
}

func Unfavorite(ctx context.Context, req *video.UnfavoriteRequest) error {
	resp, err := videoClient.Unfavorite(ctx, req)
	if err != nil {
		return err
	}
	if resp.Resp.StatusCode != 0 {
		return util.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMessage)
	}
	return nil
}
