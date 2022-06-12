package rpc

import (
	"context"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/kitex_gen/user/userservice"
	"douyin-demo-micro/util"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"time"
)

var userClient userservice.Client

func initUserRPC() {
	r, err := etcd.NewEtcdResolver([]string{util.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := userservice.NewClient(
		util.UserService,
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
	userClient = c
}

func FindUser(ctx context.Context, req *user.FindUserRequest) (*user.User, error) {
	resp, err := userClient.FindUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.Resp.StatusCode, resp.Resp.Resp.StatusMessage)
	}
	return resp.Resp.User, nil
}

func MGetUser(ctx context.Context, req *user.MGetUserRequest) ([]*user.User, error) {
	resp, err := userClient.MGetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.Resp.StatusCode, resp.Resp.Resp.StatusMessage)
	}
	return resp.Resp.Users, nil
}
