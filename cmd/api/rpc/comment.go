package rpc

import (
	"context"
	"douyin-demo-micro/kitex_gen/comment"
	"douyin-demo-micro/kitex_gen/comment/commentservice"
	"douyin-demo-micro/util"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"time"
)

var commentClient commentservice.Client

func initCommentRPC() {
	r, err := etcd.NewEtcdResolver([]string{util.EtcdAddress})
	if err != nil {
		panic(err)
	}

	c, err := commentservice.NewClient(
		util.CommentService,
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
	commentClient = c
}

func CommentCreate(ctx context.Context, req *comment.CreateCommentRequest) (*comment.Comment, error) {
	resp, err := commentClient.CreateComment(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMessage)
	}
	return resp.Comment, nil
}

func CommentDelete(ctx context.Context, req *comment.DeleteCommentRequest) error {
	resp, err := commentClient.DeleteComment(ctx, req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 0 {
		return util.NewErrNo(resp.StatusCode, resp.StatusMessage)
	}
	return nil
}

func CommentList(ctx context.Context, req *comment.CommentListRequest) ([]*comment.Comment, error) {
	resp, err := commentClient.CommentList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.Resp.StatusCode != 0 {
		return nil, util.NewErrNo(resp.Resp.StatusCode, resp.Resp.StatusMessage)
	}
	return resp.CommentList, nil
}
