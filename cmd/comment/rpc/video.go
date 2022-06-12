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

func CheckVideo(ctx context.Context, req *video.CheckVideoRequest) (bool, error) {
	resp, err := videoClient.CheckVideo(ctx, req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 0 {
		return false, util.NewErrNo(resp.StatusCode, resp.StatusMessage)
	}
	return true, nil
}

func UpdateCommentCount(ctx context.Context, req *video.UpdateCommentCountRequest) (bool, error) {
	resp, err := videoClient.UpdateCommentCount(ctx, req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 0 {
		return false, util.NewErrNo(resp.StatusCode, resp.StatusMessage)
	}
	return true, nil
}
