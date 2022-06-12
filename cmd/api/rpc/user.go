package rpc

import (
	"context"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/kitex_gen/user/userservice"
	"douyin-demo-micro/util"
	"encoding/json"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/go-redis/redis/v8"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"strconv"
	"time"
)

var userClient userservice.Client
var redisClient *redis.Client

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
	redisClient = redis.NewClient(&redis.Options{
		Addr: util.REDIS_ADDRESS,
		//db:   0,
	})
}

func CreateUser(ctx context.Context, req *user.CreateUserRequest) error {
	resp, err := userClient.CreateUser(ctx, req)
	if err != nil {
		return err
	}
	if resp.Resp.StatusCode != 0 {
		return util.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMessage)
	}
	return nil
}

func CheckUser(ctx context.Context, req *user.CheckUserRequest) (*user.User, error) {
	resp, err := userClient.CheckUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.Resp.StatusCode, resp.Resp.Resp.StatusMessage)
	}
	go func() {
		_ = SetCacheUser(ctx, resp.Resp.User)
	}()
	return resp.Resp.User, nil
}

func FindUser(ctx context.Context, req *user.FindUserRequest) (*user.User, error) {
	user, err := GetCacheUser(ctx, int(req.Req.UserId))
	if err == nil {
		return user, nil
	}
	resp, err := userClient.FindUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.Resp.StatusCode, resp.Resp.Resp.StatusMessage)
	}
	go func() {
		_ = SetCacheUser(ctx, resp.Resp.User)
	}()
	return resp.Resp.User, nil
}

func GetCacheUser(ctx context.Context, userID int) (*user.User, error) {
	var user *user.User
	userBytes, err := redisClient.Get(ctx, util.USERHEAD+strconv.Itoa(userID)).Bytes()
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, err
	}
	return user, nil
}

func SetCacheUser(ctx context.Context, user *user.User) error {
	userStr, _ := json.Marshal(user)
	return redisClient.Set(ctx,
		util.USERHEAD+strconv.Itoa(int(user.Id)),
		userStr,
		util.EXPIRATION).Err()
}

//func MGetUser(ctx context.Context, req *user.MGetUserRequest)(U)

func Follow(ctx context.Context, req *user.FollowRequest) error {
	resp, err := userClient.Follow(ctx, req)
	if err != nil {
		return err
	}
	if resp.Resp.StatusCode != 0 {
		return util.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMessage)
	}
	return nil
}

func Unfollow(ctx context.Context, req *user.UnfollowRequest) error {
	resp, err := userClient.Unfollow(ctx, req)
	if err != nil {
		return err
	}
	if resp.Resp.StatusCode != 0 {
		return util.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMessage)
	}
	return nil
}

func FollowList(ctx context.Context, req *user.FollowListRequest) ([]*user.User, error) {
	resp, err := userClient.FollowList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.Resp.StatusCode, resp.Resp.Resp.StatusMessage)
	}
	return resp.Resp.Users, nil
}

func FollowerList(ctx context.Context, req *user.FollowerListRequest) ([]*user.User, error) {
	resp, err := userClient.FollowerList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.Resp.StatusCode, resp.Resp.Resp.StatusMessage)
	}
	return resp.Resp.Users, nil
}
