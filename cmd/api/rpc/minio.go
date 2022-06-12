package rpc

import (
	"douyin-demo-micro/kitex_gen/video"
	"douyin-demo-micro/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
	"os"
	"path/filepath"
)

var playPrefix string
var coverPrefix string
var minioClient *minio.Client

func initMinIO() {
	if util.STATIC {
		playPrefix = fmt.Sprintf("%s%s%s", util.ENDPOINT, util.RESOURCES, util.VIDEODIR)
		coverPrefix = fmt.Sprintf("%s%s%s", util.ENDPOINT, util.RESOURCES, util.COVERDIR)
	} else {
		playPrefix = fmt.Sprintf("http://%s/%s/%s", util.MINIO_ADDRESS, util.BUCKET_NAME, util.VIDEODIR)
		coverPrefix = fmt.Sprintf("http://%s/%s/%s", util.MINIO_ADDRESS, util.BUCKET_NAME, util.COVERDIR)
		if minioClient == nil {
			minioClient, _ = util.NewMINIOClient()
		}
	}
}

func resolveEndpoint(videos []*video.Video) {
	for _, v := range videos {
		v.PlayUrl = playPrefix + v.PlayUrl
		v.CoverUrl = coverPrefix + v.CoverUrl
	}
}

func UploadPlay(ctx *gin.Context, data *multipart.FileHeader, fileBase string) error {
	if !util.STATIC {
		if err := util.UploadMemFile(ctx, minioClient, util.BUCKET_NAME, data, util.VIDEODIR+fileBase); err != nil {
			return err
		}
	}
	return nil
}

func UploadCover(ctx *gin.Context, data *multipart.FileHeader, fileBase, coverBase string) error {
	playFileName := filepath.Join(util.FILEDIR, util.VIDEODIR, fileBase)
	coverFileName := filepath.Join(util.FILEDIR, util.COVERDIR, coverBase)
	// TODO 获取封面的权宜之计，但是好像都是权宜之计
	if err := ctx.SaveUploadedFile(data, playFileName); err != nil {
		return nil
	}
	// 取第一帧作为封面
	if err := util.ExtractCover(playFileName, coverFileName); err != nil {
		return nil
	}
	if !util.STATIC {
		if err := util.Upload(ctx, minioClient, util.BUCKET_NAME, coverFileName, util.COVERDIR+coverBase); err != nil {
			return err
		}
		_ = os.Remove(playFileName)
		_ = os.Remove(coverFileName)
	}
	return nil
}
