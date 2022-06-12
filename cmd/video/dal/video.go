package dal

import (
	"context"
	"douyin-demo-micro/util"
	"gorm.io/gorm"
	"time"
)

// Video 加外键到User可能会导致查询变慢，但是可以减少查询次数，需要考虑
// 一致性检查放到service层
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
	return util.VideoTable
}

func CreateVideo(ctx context.Context, video *Video) error {
	err := DB.WithContext(ctx).Create(video).Error
	if err != nil {
		return util.ErrInternalError
	}
	return nil
}

func Feed(ctx context.Context, latestTime *int64) ([]*Video, error) {
	tmpDB := DB
	if latestTime != nil && util.ValidTime(*latestTime) {
		timestamp := time.UnixMilli(*latestTime).Format("2006-01-02 15:04:05")
		tmpDB = tmpDB.Where("updated_at < ?", timestamp)
	}
	var videos []*Video
	if err := tmpDB.WithContext(ctx).Order("created_at DESC").Limit(30).Find(&videos).Error; err != nil {
		return nil, util.ErrInternalError
	}
	return videos, nil
}

func QueryVideosByAuthorID(ctx context.Context, userID int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	if err := DB.WithContext(ctx).Where("author_id=?", userID).Limit(100).Find(&videos).Error; err != nil {
		return nil, util.ErrInternalError
	}
	return videos, nil
}

func QueryVideosByIDs(ctx context.Context, videoIDs []int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	if err := DB.WithContext(ctx).Find(&videos, videoIDs).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func CheckVideo(ctx context.Context, videoID int64) error {
	var video Video
	if err := DB.WithContext(ctx).Find(&video, videoID).Error; err != nil || video.ID == 0 {
		return util.ErrIllegalArguments
	}
	return nil
}

func UpdateCommentCount(ctx context.Context, videoID, delta int64) error {
	if err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var video Video
		if err := tx.Find(&video, videoID).Error; err != nil || videoID == 0 {
			return util.ErrInternalError
		}
		if err := tx.Model(&video).Update("comment_count", util.MaxI64(video.CommentCount+delta, 0)).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return util.ErrInternalError
	}
	return nil
}

func UpdateFavoriteCount(ctx context.Context, videoID, delta int64) error {
	if err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 目前没有允许删除视频，应该不需要考虑查不到视频的情况
		var video Video
		if err := tx.Find(&video, videoID).Error; err != nil {
			return util.ErrInternalError
		}
		// 当使用了Model方法，且该对象主键有值，该值会被用于构建条件。即"WHERE id=?"
		if err := tx.Model(&video).Update("favorite_count", util.MaxI64(video.FavoriteCount+delta, 0)).Error; err != nil {
			return util.ErrInternalError
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
