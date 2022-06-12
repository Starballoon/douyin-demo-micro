package util

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
)

// ExtractCover 取第一帧作为封面
func ExtractCover(playFileName, coverFileName string) error {
	cmd := exec.CommandContext(context.Background(), FFMPEGEXE,
		//"-loglevel", "error",
		// 覆盖输出文件
		"-y",
		// 偏移量，ffmpeg似乎第一帧就是从1开始，而不是0
		"-ss", "1",
		// 持续时长
		"-t", "1",
		"-i", playFileName,
		// 输出帧数
		"-vframes", "1",
		coverFileName)
	return cmd.Run()
}

// NewMINIOClient 构造OSS的客户端
func NewMINIOClient() (*minio.Client, error) {
	return minio.New(MINIO_ADDRESS, &minio.Options{
		Creds:  credentials.NewStaticV4(ACCESS_KEY_ID, ACCESS_ACCESS_KEY, ""),
		Secure: false,
	})
}

// Upload 重复上传不会报错
func Upload(ctx context.Context, client *minio.Client, bucket, localFile, remoteFile string) error {
	info, err := client.StatObject(ctx, bucket, remoteFile, minio.StatObjectOptions{})
	if len(info.Key) != 0 && info.Size != 0 {
		return nil
	}
	file, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, _ := file.Stat()
	_, err = client.PutObject(context.Background(), bucket, remoteFile, file, stat.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

func UploadMemFile(ctx context.Context, client *minio.Client, bucket string, localFile *multipart.FileHeader, remoteFile string) error {
	info, err := client.StatObject(ctx, bucket, remoteFile, minio.StatObjectOptions{})
	if len(info.Key) != 0 && info.Size != 0 {
		return nil
	}
	file, err := localFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = client.PutObject(context.Background(), bucket, remoteFile, file, localFile.Size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

func MD5FileHeader(test *multipart.FileHeader) (string, error) {
	file, err := test.Open()
	if err != nil {
		return "", err
	}
	md := md5.New()
	io.Copy(md, file)
	return hex.EncodeToString(md.Sum(nil)), nil
}

func MD5Filename(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	md := md5.New()
	io.Copy(md, file)
	return hex.EncodeToString(md.Sum(nil)), nil
}

//filepath.Base(strings.TrimSuffix(fileBase, path.Ext(fileBase))
