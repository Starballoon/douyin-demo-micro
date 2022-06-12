package main

import (
	"context"
	"douyin-demo-micro/cmd/comment/service"
	"douyin-demo-micro/kitex_gen/comment"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/util"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CreateComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *comment.CreateCommentRequest) (resp *comment.CreateCommentResponse, err error) {
	resp = new(comment.CreateCommentResponse)
	resp.Resp = new(user.BaseResp)

	comment, err := service.CreateComment(ctx, req)

	if err != nil {
		resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.Comment = comment
	resp.Resp = util.SuccessResp()
	return resp, nil
}

// DeleteComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) DeleteComment(ctx context.Context, req *comment.DeleteCommentRequest) (resp *user.BaseResp, err error) {
	resp = new(user.BaseResp)

	err = service.DeleteComment(ctx, req)
	if err != nil {
		resp = util.FailResp(1, err)
		return resp, nil
	}
	resp = util.SuccessResp()
	return resp, nil
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	resp = new(comment.CommentListResponse)
	resp.Resp = new(user.BaseResp)

	comments, err := service.CommentList(ctx, req)
	if err != nil {
		resp.Resp = util.FailResp(1, err)
		return resp, nil
	}
	resp.CommentList = comments
	resp.Resp = util.SuccessResp()
	return resp, nil
}
