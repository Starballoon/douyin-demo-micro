package service

import (
	"context"
	"douyin-demo-micro/cmd/comment/dal"
	"douyin-demo-micro/cmd/comment/rpc"
	"douyin-demo-micro/kitex_gen/comment"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/kitex_gen/video"
	"douyin-demo-micro/util"
)

// CreateComment TODO 没有最终一致性
func CreateComment(ctx context.Context, req *comment.CreateCommentRequest) (*comment.Comment, error) {
	if exists, err := rpc.CheckVideo(ctx, &video.CheckVideoRequest{VideoId: req.VideoId}); err != nil || !exists {
		return nil, err
	}
	user, err := rpc.FindUser(ctx, &user.FindUserRequest{Req: &user.IdRequest{UserId: req.ReviewerId}})
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, util.ErrIllegalArguments
	}
	comment1 := &dal.Comment{
		VideoID:    req.VideoId,
		ReviewerID: req.ReviewerId,
		Content:    req.Content,
	}
	err = dal.CreateComment(ctx, comment1)
	if err != nil {
		return nil, err
	}
	ok, err := rpc.UpdateCommentCount(ctx, &video.UpdateCommentCountRequest{
		VideoId: comment1.VideoID,
		Delta:   1,
	})
	if err != nil || !ok {
		_ = dal.DeleteComment(ctx, comment1)
		return nil, err
	}
	return Comment(comment1, user), nil
}

func DeleteComment(ctx context.Context, req *comment.DeleteCommentRequest) error {
	if req.CommentId == 0 || req.UserId == 0 {
		return util.ErrIllegalArguments
	}
	// 无事务保护，目前只允许删除自己评论，并发删除自己评论的场景应该不存在
	comment1, err := dal.FindComment(ctx, req.CommentId)
	if err != nil || comment1.ReviewerID != req.UserId {
		return util.ErrIllegalArguments
	}
	err = dal.DeleteComment(ctx, comment1)
	if err != nil {
		return util.ErrInternalError
	}
	ok, err := rpc.UpdateCommentCount(ctx, &video.UpdateCommentCountRequest{
		VideoId: comment1.VideoID,
		Delta:   -1,
	})
	if err != nil || !ok {
		return util.ErrInternalError
	}
	return nil
}

func CommentList(ctx context.Context, req *comment.CommentListRequest) ([]*comment.Comment, error) {
	if req.VideoId == 0 {
		return nil, util.ErrIllegalArguments
	}
	comments, err := dal.CommentListByVideoID(ctx, req.VideoId)
	if err != nil {
		return nil, util.ErrInternalError
	}
	if len(comments) == 0 {
		return []*comment.Comment{}, nil
	}
	// user ID 去重
	reviewerIDs := dedupID(comments)

	reviewers, err := rpc.MGetUser(ctx, &user.MGetUserRequest{
		UserIds: reviewerIDs,
	})
	if err != nil {
		return nil, util.ErrInternalError
	}

	results := Comments(comments, reviewers)
	return results, nil
}

// dedupID comment reviewerID 去重
func dedupID(comments []*dal.Comment) []int64 {
	reviewerIDMap := make(map[int64]struct{})
	for _, c := range comments {
		reviewerIDMap[c.ReviewerID] = struct{}{}
	}
	reviewerIDs := make([]int64, 0, len(reviewerIDMap))
	for reviewerID := range reviewerIDMap {
		reviewerIDs = append(reviewerIDs, reviewerID)
	}
	return reviewerIDs
}
