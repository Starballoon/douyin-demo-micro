package handlers

import (
	"douyin-demo-micro/cmd/api/rpc"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/kitex_gen/video"
	"douyin-demo-micro/util"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(ctx *gin.Context) {
	user1, err := findUser(ctx)
	if err != nil {
		failResp(ctx, 1, err)
	}
	action, err := strconv.Atoi(ctx.Query("action_type"))
	// 1-点赞， 2-取消点赞
	if err != nil || !(action == 1 || action == 2) {
		failResp(ctx, 1, util.ErrIllegalArguments)
		return
	}
	videoID, err := strconv.Atoi(ctx.Query("video_id"))
	if err != nil {
		failResp(ctx, 1, util.ErrIllegalArguments)
		return
	}
	switch action {
	case util.CREATE_FAVORITE:
		err = rpc.Favorite(ctx, &video.FavoriteRequest{
			VideoId: int64(videoID),
			UserId:  user1.Id,
		})
	case util.DELETE_FAVORITE:
		err = rpc.Unfavorite(ctx, &video.UnfavoriteRequest{
			VideoId: int64(videoID),
			UserId:  user1.Id,
		})
	}
	if err != nil {
		failResp(ctx, 1, err)
		return
	}
	successResp(ctx)
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	// 如果token中的id和request body中的user_id字段不同，该表现出什么行为？
	// 刷视频的时候它会拉取当前视频作者的信息，但是为什么拉当前作者的喜欢列表？？？
	user1, err := findUser(c)
	if err != nil || user1 == nil {
		failResp(c, 1, jwt.ErrFailedAuthentication)
	}

	videos, err := rpc.FavoriteList(c, &video.FavoriteListRequest{Req: &user.IdRequest{UserId: user1.Id}})
	if err != nil {
		failResp(c, 1, util.ErrInternalError)
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		//VideoList: DemoVideos,
		VideoList: videos,
	})
}
