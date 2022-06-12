package handlers

import (
	"douyin-demo-micro/cmd/api/rpc"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/kitex_gen/video"
	"douyin-demo-micro/util"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VideoListResponse struct {
	Response
	VideoList []*video.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(ctx *gin.Context) {
	var err error
	user, err := findUser(ctx)
	if err != nil {
		failResp(ctx, 1, jwt.ErrFailedAuthentication)
	}
	title := ctx.PostForm("title")
	data, err := ctx.FormFile("data")
	if len(title) == 0 || err != nil || data == nil {
		failResp(ctx, 1, util.ErrIllegalArguments)
	}

	// 视频和封面存成同名，名字为视频MD5摘要值
	nameBase, _ := util.MD5FileHeader(data)
	fileBase := nameBase + ".mp4"
	coverBase := nameBase + ".jpg"

	done := make(chan struct{})
	go func(err error) {
		if e := rpc.PublishVideo(ctx, &video.CreateVideoRequest{
			AuthorId:      user.Id,
			VideoFilename: fileBase,
			CoverFilename: coverBase,
			Title:         title,
		}); e != nil {
			err = e
		}
		done <- struct{}{}
	}(err)

	go func(err error) {
		if e := rpc.UploadPlay(ctx, data, fileBase); e != nil {
			err = e
		}
		done <- struct{}{}
	}(err)

	go func(err error) {
		if e := rpc.UploadCover(ctx, data, fileBase, coverBase); e != nil {
			err = e
		}
		done <- struct{}{}
	}(err)

	for i := 0; i < 3; i += 1 {
		<-done
	}
	if err != nil {
		failResp(ctx, 1, err)
		return
	}
	successResp(ctx)
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	// 如果token中的id和request body中的user_id字段不同，该表现出什么行为？
	// 刷视频的时候它会拉取当前视频作者的信息
	user1, err := findUser(c)
	if err != nil {
		failResp(c, 1, jwt.ErrFailedAuthentication)
	}
	videos, err := rpc.PublishList(c, &video.PublishListRequest{Req: &user.IdRequest{UserId: user1.Id}})
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
