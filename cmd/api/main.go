package main

import (
	"context"
	"douyin-demo-micro/cmd/api/rpc"
	"douyin-demo-micro/util"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// InitMINIO 检查并初始化OSS存储，仅检查一次
func InitMINIO() error {
	minioClient, err := util.NewMINIOClient()
	if err != nil {
		return err
	}
	exists, err := minioClient.BucketExists(context.Background(), util.BUCKET_NAME)
	if err != nil {
		return err
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), util.BUCKET_NAME, minio.MakeBucketOptions{
			Region:        "",
			ObjectLocking: false,
		})
		if err != nil {
			return err
		}
		err = minioClient.SetBucketPolicy(context.Background(), util.BUCKET_NAME, util.POLICY)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// TODO 伪造数据 注意删除
	util.FakeData()

	util.InitJaeger(util.ApiService)

	rpc.InitRPC()

	err := InitMINIO()
	if err != nil {
		panic(err)
	}

	// TODO
	rpc.InitRPC()

	r := gin.Default()
	initRouter(r)

	if err = r.Run(util.ApiServicePort); err != nil {
		klog.Fatal(err)
	}
}
