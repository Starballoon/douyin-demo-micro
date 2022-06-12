package util

import "time"

const (
	// ADDRESS 注意数据库、网关和OSS存储地址可以不同
	ADDRESS = "192.168.50.144"
	PORT    = ":8080"
	// ENDPOINT 下面ip和端口是静态资源的访问点，OSS存储需要和服务分开启动
	// ENDPOINT     = "http://localhost:8080/"
	ENDPOINT  = "http://" + ADDRESS + PORT + "/"
	RESOURCES = "static/"
	VIDEODIR  = "videos/"
	COVERDIR  = "covers/"
	//FILEDIR      = "./public"
	FILEDIR      = "/home/zxq/tmp/douyin-demo-micro/public/"
	DEFAULTCOVER = "logo.png"

	// 数据库相关
	ACCOUNT  = "root"
	PASSWORD = "root"
	DB_PORT  = "3306"
	DATABASE = "douyin"
	// DSN    = "root:LLQtT$3v@tcp(192.168.50.144:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	DSN           = ACCOUNT + ":" + PASSWORD + "@tcp(" + ADDRESS + ":" + DB_PORT + ")/" + DATABASE + "?charset=utf8mb4&parseTime=True&loc=Local"
	UserTable     = "user"
	VideoTable    = "video"
	CommentTable  = "comment"
	FollowTable   = "follow"
	FavoriteTable = "favorite"

	// STATIC 静态资源或MINIO资源
	STATIC = false

	IDENTITY = "id"
	JWTKEY   = "freshbird"

	// MINIO OSS存储
	//MINIO_ADDRESS = "localhost:9000"
	MINIO_ADDRESS     = ADDRESS + ":9000"
	ACCESS_KEY_ID     = "minioadmin"
	ACCESS_ACCESS_KEY = "minioadmin"
	BUCKET_NAME       = "douyin"
	// POLICY MINIO 存储桶的访问策略，大意是允许 任何人 获取 douyin桶内的文件
	POLICY = "{\n\"Version\":\"2012-10-17\",\n\"Statement\":[\n{\n\"Effect\":\"Allow\",\n\"Principal\":{\n\"AWS\":[\n\"*\"\n]\n},\n\"Action\":[\n\"s3:GetBucketLocation\",\n\"s3:ListBucket\",\n\"s3:ListBucketMultipartUploads\"\n],\n\"Resource\":[\n\"arn:aws:s3:::" +
		BUCKET_NAME + "\"\n]\n},\n{\n\"Effect\":\"Allow\",\n\"Principal\":{\n\"AWS\":[\n\"*\"\n]\n},\n\"Action\":[\n\"s3:DeleteObject\",\n\"s3:GetObject\",\n\"s3:ListMultipartUploadParts\",\n\"s3:PutObject\",\n\"s3:AbortMultipartUpload\"\n],\n\"Resource\":[\n\"arn:aws:s3:::douyin/*\"\n]\n}\n]\n}"

	// REDIS_ADDRESS Redis用作缓存，目前是无加密操作
	REDIS_ADDRESS = ADDRESS + ":6379"
	USERHEAD      = "DOUYIN_USERID_"
	// EXPIRATION 缓存过期时间
	EXPIRATION = time.Minute * 30

	// 关系操作模式常量
	// 1 - create, 2 - delete
	CREATE_FAVORITE = 1
	DELETE_FAVORITE = 2
	CREATE_COMMENT  = 1
	DELETE_COMMENT  = 2
	CREATE_FOLLOW   = 1
	DELETE_FOLLOW   = 2

	// FFMPEGEXE 服务器ffmpeg可执行文件绝对路径
	//FFMPEGEXE = "C:/envs/ffmpeg/bin/ffmpeg.exe"
	FFMPEGEXE = "/usr/bin/ffmpeg"

	JWT_TMP_KEY = "JWT_Temp_Valid"

	EtcdAddress                = ADDRESS + ":2379"
	ApiService                 = "douyin-api"
	UserService                = "douyin-user"
	VideoService               = "douyin-video"
	CommentService             = "douyin-comment"
	ApiServicePort             = ":8080"
	UserServicePort            = ":5714"
	VideoServicePort           = ":5715"
	CommentServicePort         = ":5716"
	CPURateLimit       float64 = 80.0
)
