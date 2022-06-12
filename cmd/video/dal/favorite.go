package dal

import (
	"context"
	"douyin-demo-micro/util"
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	VideoID int64 `gorm:"column:video_id;index:positive;index:reverse"`
	UserID  int64 `gorm:"column:user_id;index:positive;index:reverse,priority:9"`
}

func (u *Favorite) TableName() string {
	return util.FavoriteTable
}

func QueryUnscoped(ctx context.Context, videoID, userID int64) (*Favorite, error) {
	// TODO 似乎抖音允许自己点赞自己视频，不允许点赞自己评论？
	var favorite Favorite
	if err := DB.WithContext(ctx).Unscoped().
		Where("video_id=? AND user_id=?", videoID, userID).Find(&favorite).Error; err != nil {
		return nil, util.ErrInternalError
	}
	return &favorite, nil
}

func CreateFavor(ctx context.Context, favorite *Favorite) error {
	if err := DB.WithContext(ctx).Create(favorite).Error; err != nil {
		return util.ErrInternalError
	}
	return nil
}

func DeleteFavor(ctx context.Context, favorite *Favorite) error {
	if err := DB.WithContext(ctx).Delete(&favorite).Error; err != nil {
		return util.ErrInternalError
	}
	return nil
}

func RecoverFavor(ctx context.Context, favorite *Favorite) error {
	if err := DB.WithContext(ctx).Unscoped().Model(&favorite).
		Update("deleted_at", nil).Error; err != nil {
		return util.ErrInternalError
	}
	return nil
}

func FavoritesByIDs(ctx context.Context, videoIDs []int64, userID int64) ([]int64, error) {
	result := make([]int64, 0)
	if err := DB.WithContext(ctx).Table(util.FavoriteTable).Select("video_id").
		Where("video_id IN ? AND user_id=?", videoIDs, userID).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func FavoritesByID(ctx context.Context, userID int64) ([]int64, error) {
	videoIDs := make([]int64, 0)
	if err := DB.WithContext(ctx).Table(util.FavoriteTable).Select("video_id").
		Where("user_id=?", userID).Limit(100).Find(&videoIDs).Error; err != nil {
		return nil, err
	}
	return videoIDs, nil
}
