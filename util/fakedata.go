package util

import (
	"context"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

var names = []string{
	"陈玟玉", "陈和煦", "陈南风", "陈春海",
	"陈云亭", "陈志", "陈思雅", "陈浩邈",
	"陈山菡", "陈白安", "陈嘉宝", "窦又琴",
	"窦仁", "窦才艺", "窦碧春", "窦从霜",
	"黄丽华", "黄妙之", "黄运珹", "黄娜娜",
	"黄烨然", "李秋灵", "李迎曼", "李听云",
	"李凝安", "李倚", "邹实", "邹永言",
	"邹彭越", "邹惜筠", "邹子楠", "赵朔",
	"赵坚", "赵念珍", "赵侠", "赵又槐"}

var db *gorm.DB

type User struct {
	gorm.Model
	EmailAccount  string  `gorm:"column:email_account;index:unique_user,priority:8"`
	EmailDomain   string  `gorm:"column:email_domain;index:unique_user,priority:9"`
	Password      *[]byte `gorm:"column:password;type:BINARY(64)"`
	Name          string  `gorm:"column:name"`
	FollowCount   int64   `gorm:"column:follow_count"`
	FollowerCount int64   `gorm:"column:follower_count"`
}

func (u *User) TableName() string {
	return UserTable
}

type Video struct {
	gorm.Model
	AuthorID      int64  `gorm:"column:author_id;index"`
	VideoFilename string `gorm:"column:video_filename"`
	CoverFilename string `gorm:"column:cover_filename"`
	FavoriteCount int64  `gorm:"column:favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count"`
	Title         string `gorm:"column:title;index"`
}

func (u *Video) TableName() string {
	return VideoTable
}

type Comment struct {
	gorm.Model
	VideoID    int64  `gorm:"column:video_id;index"`
	ReviewerID int64  `gorm:"column:reviewer_id;index"`
	Content    string `gorm:"column:content"`
}

func (u *Comment) TableName() string {
	return CommentTable
}

type Following struct {
	gorm.Model
	LeaderID   int64 `gorm:"column:leader_id;index:positive;index:reverse"`
	FollowerID int64 `gorm:"column:follower_id;index:positive;index:reverse,priority:9"`
}

func (u *Following) TableName() string {
	return FollowTable
}

type Favorite struct {
	gorm.Model
	VideoID int64 `gorm:"column:video_id;index:positive;index:reverse"`
	UserID  int64 `gorm:"column:user_id;index:positive;index:reverse,priority:9"`
}

func (u *Favorite) TableName() string {
	return FavoriteTable
}
func InitDB() error {
	if db != nil {
		return errors.New("multiple initialization is not allowed")
	}
	var err error
	db, err = gorm.Open(mysql.Open(DSN),
		&gorm.Config{
			PrepareStmt: true,
		})
	if err != nil {
		return err
	}
	return InitTables()
}

func InitTables() error {
	var err error
	m := db.Migrator()
	if !m.HasTable(&User{}) {
		err = m.CreateTable(&User{})
		if err != nil {
			return err
		}
	}
	if !m.HasTable(&Video{}) {
		err = m.CreateTable(&Video{})
		if err != nil {
			return err
		}
	}
	if !m.HasTable(&Comment{}) {
		err = m.CreateTable(&Comment{})
		if err != nil {
			return err
		}
	}
	if !m.HasTable(&Following{}) {
		err = m.CreateTable(&Following{})
		if err != nil {
			return err
		}
	}
	if !m.HasTable(&Favorite{}) {
		err = m.CreateTable(&Favorite{})
		if err != nil {
			return err
		}
	}
	return err
}

func FakeData() error {
	InitDB()
	var count int64
	db.WithContext(context.Background()).Model(&Video{}).Count(&count)
	if count >= int64(len(names)) {
		return nil
	}
	err := fakeUsers()
	if err != nil {
		return err
	}
	err = fakeVideos()
	if err != nil {
		return err
	}
	err = fakeFollows()
	if err != nil {
		return err
	}
	err = fakeFavorite()
	if err != nil {
		return nil
	}
	err = fakeComments()
	if err != nil {
		return nil
	}
	return nil
}

func fakeUsers() error {
	// 每个用户发布1个视频
	for i := len(names) - 1; i > 0; i -= 1 {
		j := rand.Intn(i + 1)
		names[i], names[j] = names[j], names[i]
	}
	passwd := Encrypt("12345678")
	user := &User{
		EmailDomain:   "bytedance.com",
		Password:      &passwd,
		FollowCount:   1,
		FollowerCount: 1,
	}
	for i := 1; i < len(names); i += 1 {
		user.Model = gorm.Model{}
		user.Name = names[i-1]
		user.EmailAccount = "user" + strconv.Itoa(i)
		err := db.WithContext(context.Background()).Create(user).Error
		if err != nil {
			return err
		}
	}
	user.Model = gorm.Model{}
	user.Name = names[len(names)-1]
	user.FollowerCount = int64(len(names) - 1)
	user.FollowCount = int64(len(names) - 1)
	user.EmailAccount = "user" + strconv.Itoa(len(names))
	err := db.WithContext(context.Background()).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func fakeVideos() error {
	// 一共36个视频
	video := &Video{
		CommentCount:  1,
		FavoriteCount: 1,
	}

	minioClient, _ := NewMINIOClient()

	for i := 1; i < len(names); i += 1 {
		video.Model = gorm.Model{}
		video.AuthorID = int64(i)
		video.VideoFilename = "video" + strconv.Itoa(i) + ".mp4"
		video.CoverFilename = "cover" + strconv.Itoa(i) + ".jpg"
		video.Title = "video" + strconv.Itoa(i) + ".mp4"
		video.UpdatedAt = time.Now().Add(time.Second * 2 * time.Duration(i))
		// 手动提取封面
		//ExtractCover(
		//	filepath.Join(FILEDIR, VIDEODIR, "video"+strconv.Itoa(i)+".mp4"),
		//	filepath.Join(FILEDIR, COVERDIR, "cover"+strconv.Itoa(i)+".jpg"))
		_ = Upload(context.Background(), minioClient, BUCKET_NAME,
			FILEDIR+VIDEODIR+video.VideoFilename, VIDEODIR+video.VideoFilename)
		_ = Upload(context.Background(), minioClient, BUCKET_NAME,
			FILEDIR+COVERDIR+video.CoverFilename, COVERDIR+video.CoverFilename)
		err := db.WithContext(context.Background()).Create(video).Error
		if err != nil {
			return err
		}
	}

	video.Model = gorm.Model{}
	video.AuthorID = int64(len(names))
	video.VideoFilename = "video" + strconv.Itoa(len(names)) + ".mp4"
	video.CoverFilename = "cover" + strconv.Itoa(len(names)) + ".jpg"
	video.Title = "video" + strconv.Itoa(len(names)) + ".mp4"
	video.FavoriteCount = int64(len(names))
	video.CommentCount = int64(len(names) + 1)
	//ExtractCover(
	//	filepath.Join(FILEDIR, VIDEODIR, video.VideoFilename),
	//	filepath.Join(FILEDIR, COVERDIR, video.CoverFilename))
	_ = Upload(context.Background(), minioClient, BUCKET_NAME,
		FILEDIR+VIDEODIR+video.VideoFilename, VIDEODIR+video.VideoFilename)
	_ = Upload(context.Background(), minioClient, BUCKET_NAME,
		FILEDIR+COVERDIR+video.CoverFilename, COVERDIR+video.CoverFilename)
	err := db.WithContext(context.Background()).Create(video).Error
	if err != nil {
		return err
	}
	return nil
}

func fakeFollows() error {
	// 36号用户关注了所有人
	for i := 1; i < len(names); i += 1 {
		err := db.WithContext(context.Background()).Create(&Following{
			LeaderID:   int64(i),
			FollowerID: int64(len(names)),
		}).Error
		if err != nil {
			return nil
		}
	}
	// 所有用户关注36号用户
	for i := 1; i < len(names); i += 1 {
		err := db.WithContext(context.Background()).Create(&Following{
			LeaderID:   int64(len(names)),
			FollowerID: int64(i),
		}).Error
		if err != nil {
			return nil
		}
	}
	return nil
}

func fakeFavorite() error {
	// 36号用户喜欢所有视频
	for i := 1; i <= len(names); i += 1 {
		err := db.WithContext(context.Background()).Create(&Favorite{
			VideoID: int64(i),
			UserID:  int64(len(names)),
		}).Error
		if err != nil {
			return nil
		}
	}
	// 所有用户喜欢36号视频
	for i := 1; i <= len(names); i += 1 {
		err := db.WithContext(context.Background()).Create(&Favorite{
			VideoID: int64(len(names)),
			UserID:  int64(i),
		}).Error
		if err != nil {
			return nil
		}
	}
	return nil
}

func fakeComments() error {
	// 36号用户评论所有视频
	for i := 1; i <= len(names); i += 1 {
		err := db.WithContext(context.Background()).Create(&Comment{
			VideoID:    int64(i),
			ReviewerID: int64(len(names)),
			Content:    "This is video" + strconv.Itoa(i),
		}).Error
		if err != nil {
			return nil
		}
	}
	// 所有用户评论36号视频
	for i := 1; i <= len(names); i += 1 {
		err := db.WithContext(context.Background()).Create(&Comment{
			VideoID:    int64(len(names)),
			ReviewerID: int64(i),
			Content:    "This is user" + strconv.Itoa(i),
		}).Error
		if err != nil {
			return nil
		}
	}
	return nil
}
