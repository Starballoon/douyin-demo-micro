package handlers

import (
	"douyin-demo-micro/cmd/api/rpc"
	"douyin-demo-micro/kitex_gen/comment"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentResponse struct {
	Response
	Comment *comment.Comment `json:"comment, omitempty"`
}

type CommentListResponse struct {
	Response
	CommentList []*comment.Comment `json:"comment_list,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(ctx *gin.Context) {
	user, err := findUser(ctx)
	if err != nil {
		failResp(ctx, 1, util.ErrIllegalArguments)
	}
	// 1-发布评论，2-删除评论
	action, err := strconv.Atoi(ctx.Query("action_type"))
	if err != nil || !(action == 1 || action == 2) {
		failResp(ctx, 1, util.ErrIllegalArguments)
		return
	}
	switch action {
	case util.DELETE_COMMENT:
		commentDelete(ctx, user)
	case util.CREATE_COMMENT:
		commentCreate(ctx, user)
	}
}

func commentCreate(ctx *gin.Context, user *user.User) {
	content := ""
	content = ctx.Query("comment_text")
	if len(content) == 0 {
		failResp(ctx, 1, util.ErrIllegalArguments)
		return
	}
	videoID, err := strconv.Atoi(ctx.Query("video_id"))
	if err != nil || videoID < 1 {
		failResp(ctx, 1, util.ErrIllegalArguments)
		return
	}

	comment, err := rpc.CommentCreate(ctx, &comment.CreateCommentRequest{
		ReviewerId: user.Id,
		VideoId:    int64(videoID),
		Content:    content,
	})
	if err != nil || comment == nil {
		failResp(ctx, 1, util.ErrInternalError)
	}
	ctx.JSON(http.StatusOK, &CommentResponse{
		Response: Response{StatusCode: 0},
		Comment:  comment,
	})
}

func commentDelete(ctx *gin.Context, user *user.User) {
	commentID := 0
	// 前端限制了只能删自己的
	commentID, err := strconv.Atoi(ctx.Query("comment_id"))
	if err != nil {
		failResp(ctx, 1, util.ErrIllegalArguments)
		return
	}
	err = rpc.CommentDelete(ctx, &comment.DeleteCommentRequest{
		UserId:    user.Id,
		CommentId: int64(commentID),
	})
	if err != nil {
		failResp(ctx, 1, util.ErrInternalError)
		return
	}
	successResp(ctx)
	return
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoID, err := strconv.Atoi(c.Query("video_id"))
	if err != nil || videoID < 1 {
		failResp(c, 1, util.ErrIllegalArguments)
		return
	}
	comments, err := rpc.CommentList(c, &comment.CommentListRequest{VideoId: int64(videoID)})
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
