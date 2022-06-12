package handlers

import (
	"douyin-demo-micro/cmd/api/rpc"
	"douyin-demo-micro/kitex_gen/video"
	"douyin-demo-micro/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FeedResponse struct {
	Response
	VideoList []*video.Video `json:"video_list,omitempty"`
	NextTime  *int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(ctx *gin.Context) {
	userID, _ := identify(ctx)
	latestTime := ctx.Query("latest_time")
	req := &video.FeedRequest{}
	if userID > 0 {
		req.UserId = &userID
	}
	if len(latestTime) > 0 {
		tmp, err := strconv.ParseInt(latestTime, 0, 64)
		if err == nil {
			req.LatestTime = &tmp
		}
	}

	videos, nextTime, err := rpc.Feed(ctx, req)
	if err != nil {
		failResp(ctx, 1, util.ErrInternalError)
		return
	}
	ctx.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		//NextTime:  time.Now().UnixMicro(),
		NextTime: nextTime,
	})
}
